package giglet

import (
	"bytes"
	"crypto/tls"
	"errors"
	"giglet/safe"
	"giglet/specs"
	"io"
	"net"
	"runtime"
	"time"
)


func (server *Server) Serve(listener net.Listener) error {
	if listener == nil {
		return errors.New("empty listener")
	} else if server.isShuttingdown.Load() {
		return ErrorServerClosed
	}

	server.listenerTrack.Add(1)
	defer server.listenerTrack.Done()

	for {
		conn, err := listener.Accept()
		if server.isShuttingdown.Load() {
			if err == nil {
				conn.Close()
			}
			return ErrorServerClosed
		} else if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				time.Sleep(time.Second)
				continue
			}
			return err
		}
		go server.handle(conn)
	}
}

var bufioReaderPool safe.BufioReaderPool

func (server *Server) handle(conn net.Conn) {
	if server.Handler == nil || server.isShuttingdown.Load() {
		conn.Close()
		return
	}
	handler := server.Handler

	if tlsConn, ok := conn.(*tls.Conn); ok {
		server.applyReadTimeout(conn)
		server.applyWriteTimeout(conn)

		if err := tlsConn.Handshake(); err != nil {
			// If the handshake failed due to the client not speaking
			// TLS, assume they're speaking plaintext HTTP and write a
			// 400 response on the TLS conn's underlying net.Conn.
			if re, ok := err.(tls.RecordHeaderError); ok && re.Conn != nil {
				re.Conn.Write(responseDowngradeHTTPS)
				re.Conn.Close()
				return
			}
			if server.Debug {
				server.logger().Printf("http: tls handshake error from %s: %v", conn.RemoteAddr(), err)
			}
			return
		}
		conn.SetDeadline(zeroTime)

		proto := tlsConn.ConnectionState().NegotiatedProtocol

		if server.nextProtos != nil {
			if handler, ok := server.nextProtos[proto]; ok {
				handler(tlsConn)
				return
			}
		}
	}

	if server.isShuttingdown.Load() {
		conn.Close()
		return
	}

	reader := bufioReaderPool.Get(conn)

	defer func() {
		if err := recover(); err != nil && err != ErrorAbortHandler { 
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]

			if server.Debug {
				server.logger().Printf("http: panic serving %v: %v\n%s", conn.RemoteAddr(), err, buf)
			}
		}

		conn.Close()
		bufioReaderPool.Put(reader)
	}()

	for {
		if server.ContentMaxSizeBytes > 0 {
			reader.Reset(io.LimitReader(conn, server.ContentMaxSizeBytes))
		} else if DefaultContentMaxSizeBytes > 0 {
			reader.Reset(io.LimitReader(conn, DefaultContentMaxSizeBytes))
		}

		server.applyReadTimeout(conn)
		req, err := readRequest(reader)
		conn.SetReadDeadline(zeroTime)

		req.server = server
		req.conn = conn
		
		if err != nil {
			if server.Debug {
				server.logger().Printf("http: read request error from %s: %v", conn.RemoteAddr(), err)
			}
			switch {
			case err == ErrorTooLarge:
				conn.Write(responseRequestHeadersTooLarge)
	
			case err == ErrorUnsupportedEncoding:
				conn.Write(responseUnsupportedEncoding)
				
			default:
				if !isCommonNetReadError(err) {
					if serr, ok := err.(*statusErrorResponse); ok {
						serr.Write(conn)
					} else {
						conn.Write(responseNotProcessableError)
					}
				}

			}
			break
		}
		
		resp := handler(req)
		var header *specs.Header
		var code specs.StatusCode
		var writable WritableResponse
		if resp != nil {
			if prep, ok := resp.(PreparableResponse); ok {
				prep.Prepare()
			}

			header = resp.Header()
			code = resp.StatusCode()
			writable, _ = resp.(WritableResponse)
		}
		if header == nil {
			header = &specs.Header{}
		}

		if len(server.ServerName) > 0 {
			header.Set("Server", server.ServerName)
		} else if len(DefaultServerName) > 0 {
			header.Set("Server", DefaultServerName)
		}

		header.Set("Date", time.Now().Format(specs.TimeFormat))

		if !code.IsValid() {
			if !req.Method().CanHaveResponseBody() || writable == nil {
				code = specs.StatusCodeNoContent
			} else {
				code = specs.StatusCodeOK
			}
		}

		_, err = WriteResponseHeadTo(conn, req.ProtoAtLeast(1, 1), code, header)
		if err != nil {
			if server.Debug {
				server.logger().Printf("http: send response head to %s error: %v", conn.RemoteAddr(), err)
			}
			break
		}
		if req.method.CanHaveResponseBody() && writable != nil {
			if server.WriteTimeout > 0 {
				server.applyWriteTimeout(conn)
			}

			writable.WriteBody(conn)
			
			if server.WriteTimeout > 0 {
				conn.SetWriteDeadline(zeroTime)
			}
		}

		if req.cachedMultipart != nil {
			req.cachedMultipart.RemoveAll()
		}

		if server.isShuttingdown.Load() {
			break
		} else if req.hijacker != nil {
			req.hijacker(conn)
			break
		} else if req.Method() != specs.HttpMethodHead && writable == nil && code.HaveBody() {
			break
		}
	}
}

func WriteResponseHeadTo(writer io.Writer, is11 bool, code specs.StatusCode, header *specs.Header) (int64, error) {
	var headbuf bytes.Buffer
	
	if !code.IsValid() {
		code = specs.StatusCodeOK
	}

	_, err := code.WriteAsHeadlineTo(&headbuf, is11)
	if err != nil {
		return -1, err
	}
	if header != nil {
		_, err = header.WriteAsResponseHeaderTo(&headbuf)
		if err != nil {
			return -1, err
		}
	}
	headbuf.Write(directCrlf)
	return headbuf.WriteTo(writer)
}
