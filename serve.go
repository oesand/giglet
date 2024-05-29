package giglet

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"errors"
	"giglet/specs"
	"io"
	"net"
	"runtime"
	"time"
)


func (server *Server) Serve(listener net.Listener) error {
	if listener == nil {
		return errors.New("empty listener")
	}

	server.listenerTrack.Add(1)
	defer server.listenerTrack.Done()

	for {
		conn, err := listener.Accept()
		if err != nil {
			if server.isShuttingdown.Load() {
				return ErrorServerClosed
			} else if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				time.Sleep(time.Second)
				continue
			}
			return err
		}
		go server.work(conn)
	}
}

func (server *Server) work(conn net.Conn) {
	if server.Handler == nil {
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
				handler(conn)
				return
			}
		}
	}

	defer func() {
		if err := recover(); err != nil && err != ErrorAbortHandler {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]

			if server.Debug {
				server.logger().Printf("http: panic serving %v: %v\n%s", conn.RemoteAddr(), err, buf)
			}
		}
	}()

	// Without sync.Pool, because (*Reader).Reset are same
	reader := bufio.NewReader(conn)

	for {
		if server.ContentMaxSizeBytes > 0 {
			reader.Reset(io.LimitReader(conn, server.ContentMaxSizeBytes))
		} else if DefaultContentMaxSizeBytes > 0 {
			reader.Reset(io.LimitReader(conn, DefaultContentMaxSizeBytes))
		}

		server.applyReadTimeout(conn)
		req, err := readRequest(reader)
		conn.SetReadDeadline(zeroTime)
		
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
				if serr, ok := err.(*statusErrorResponse); ok {
					serr.Write(conn)
				} else {
					conn.Write(responseNotProcessableError)
				}

			}
			conn.Close()
			return
		}

		resp := handler(req)
		var header *specs.Header
		var code *specs.StatusCode
		var writable WritableResponse
		if resp != nil {
			if prep, ok := resp.(PreparableResponse); ok {
				prep.Prepare()
			}

			header = resp.Header()
			code = resp.StatusCode()
			writable = resp.(WritableResponse)
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

		if code == nil {
			if !req.Method().HasBody() || writable == nil {
				code = specs.StatusCodeNoContent
			} else {
				code = specs.StatusCodeOK
			}
		}

		err = WriteResponseHeadTo(conn, req.ProtoAtLeast(1, 1), code, header)
		if err != nil {
			if server.Debug {
				server.logger().Printf("http: send response head to %s error: %v", conn.RemoteAddr(), err)
			}
			return
		}
		
		if req.method.HasBody() || writable != nil {
			if server.WriteTimeout > 0 {
				server.applyWriteTimeout(conn)
			}

			writable.Respond(conn)
			
			if server.WriteTimeout > 0 {
				conn.SetWriteDeadline(zeroTime)
			}
		}

		if req.cachedMultipart != nil {
			req.cachedMultipart.RemoveAll()
		}

		if req.hijacker != nil {
			req.hijacker(conn)
			conn.Close()
			return
		}

		// conn.Close()

		

		return
		
		// [FIXME]: ADD >> ReuseConnection

	}
}

func WriteResponseHeadTo(writer io.Writer, is11 bool, code *specs.StatusCode, header *specs.Header) error {
	headbuf := &bytes.Buffer{}

	if code == nil {
		code = specs.StatusCodeOK
	}

	err := code.WriteAsHeadlineTo(headbuf, is11)
	if err != nil {
		return err
	}
	if header != nil {
		_, err = header.WriteTo(headbuf)
		if err != nil {
			return err
		}
	}
	headbuf.Write(directCrlf)
	_, err = headbuf.WriteTo(writer)
	return err
}
