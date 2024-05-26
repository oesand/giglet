package giglet

import (
	"crypto/tls"
	"net"
	"runtime"
	"time"
)


func (server *Server) Serve(listener net.Listener) error {
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
	if tlsConn, ok := conn.(*tls.Conn); ok {
		timeout := server.handshakeTimeout()
		if timeout > 0 {
			dl := time.Now().Add(timeout)
			conn.SetDeadline(dl)
		}
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
		if timeout > 0 {
			conn.SetReadDeadline(zeroTime)
			conn.SetWriteDeadline(zeroTime)
		}

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

	reader := getBufioReader(conn)

	if server.ReadTimeout > 0 {
		conn.SetReadDeadline(time.Now().Add(server.ReadTimeout))
	}

	req, err := readRequest(reader)

	if err != nil {
		switch {
		case err == ErrorTooLarge:
			conn.Write(responseRequestHeadersTooLarge)

		case err == ErrorUnsupportedEncoding:
			conn.Write(responseUnsupportedEncoding)
		}
		conn.Close()
		readerPool.Put(reader)
		return
	}


	// if server.WriteTimeout > 0 {
	// 	conn.SetWriteDeadline(time.Now().Add(server.WriteTimeout))
	// }


}
