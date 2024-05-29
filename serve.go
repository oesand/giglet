package giglet

import (
	"bufio"
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
			server.logger().Printf("http: tls handshake error from %s: %v", conn.RemoteAddr(), err)
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

	defer func() { // [FIXME]: Add continue and hijack
		if err := recover(); err != nil && err != ErrorAbortHandler {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			server.logger().Printf("http: panic serving %v: %v\n%s", conn.RemoteAddr(), err, buf)
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

		if writeStatusLine(conn, req.ProtoAtLeast(1, 1), code) != nil || 
			header.Write(conn) != nil {
			return
		}
		if req.method.HasBody() || writable != nil {
			if server.WriteTimeout > 0 {
				server.applyWriteTimeout(conn)
			}

			writable.Write(conn)
			
			if server.WriteTimeout > 0 {
				conn.SetWriteDeadline(zeroTime)
			}
		}

		return
		
		// [FIXME]: ADD >> ReuseConnection

	}
}

func writeStatusLine(writer io.Writer, is11 bool, code *specs.StatusCode) error {
	if writer == nil || code == nil { 
		return errors.New("invalid writer or status code")
	}

	var line []byte
	if is11 {
		line = append(line, httpV11...)
	} else {
		line = append(line, httpV10...)
	}

	line = append(line, ' ')
	line = code.Append(line)
	line = append(line, directCrlf...)

	_, err := writer.Write(line)
	return err
}
