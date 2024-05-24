package giglet

import (
	"crypto/tls"
	"net"
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
			conn.SetReadDeadline(dl)
			conn.SetWriteDeadline(dl)
		}
		if err := tlsConn.Handshake(); err != nil {
			// If the handshake failed due to the client not speaking
			// TLS, assume they're speaking plaintext HTTP and write a
			// 400 response on the TLS conn's underlying net.Conn.
			// [FIXME]!!!
			// if re, ok := err.(tls.RecordHeaderError); ok && re.Conn != nil && tlsRecordHeaderLooksLikeHTTP(re.RecordHeader) {
			// 	io.WriteString(re.Conn, "HTTP/1.0 400 Bad Request\r\n\r\nClient sent an HTTP request to an HTTPS server.\n")
			// 	re.Conn.Close()
			// 	return
			// }
			server.logger().Printf("http: tls handshake error from %s: %v", conn.RemoteAddr(), err)
			return
		}
		if timeout > 0 {
			conn.SetReadDeadline(time.Time{})
			conn.SetWriteDeadline(time.Time{})
		}
		// c.tlsState = new(tls.ConnectionState)
		// *c.tlsState = tlsConn.ConnectionState()
		// if proto := c.tlsState.NegotiatedProtocol; validNextProto(proto) {
		// 	if fn := c.server.TLSNextProto[proto]; fn != nil {
		// 		h := initALPNRequest{ctx, tlsConn, serverHandler{c.server}}
		// 		// Mark freshly created HTTP/2 as active and prevent any server state hooks
		// 		// from being run on these connections. This prevents closeIdleConns from
		// 		// closing such connections. See issue https://golang.org/issue/39776.
		// 		c.setState(c.rwc, StateActive, skipHooks)
		// 		fn(c.server, tlsConn, h)
		// 	}
		// 	return
		// }
	}

	// reader := getBufloReader(conn)

	// request := HttpRequest{
	// 	conn: conn,
	// }

	// for {

	// 	reader.ReadLine()

	// 	buffer, _ := reader.Peek(1)
	// 	if len(buffer) == 0 {
	// 		conn.Close()
	// 	}

	// }
}
