package upgrading

import (
	"giglet"
	"giglet/specs"
	"io"
	"net"
	"time"
)

func UpgradeProxy(req giglet.Request) giglet.Response {
	if req.ProtoNoHigher(1, 1) {
		return (&giglet.TextResponse{
			Text: "proxy: available only for httpv1",
		}).SetStatusCode(specs.StatusCodeBadGateway)
	} else if req.Method() != specs.HttpMethodConnect {
		return (&giglet.HeaderResponse{}).SetStatusCode(specs.StatusCodeMethodNotAllowed)
	}

	host := req.Header().Get("Host")
	if len(host) == 0 &&
		host != "localhost" && host != "127.0.0.1" &&
		host != "192.168.0.1" && host != "172.0.0.1" {

		return (&giglet.TextResponse{
			Text: "proxy: 'Host' header invalid, empty or not available",
		}).SetStatusCode(specs.StatusCodeBadGateway)
	}

	dest_conn, err := net.DialTimeout("tcp", host, 10 * time.Second)
	if err != nil {
		return (&giglet.TextResponse{
			Text: "proxy: destination 'Host' not available for connection",
		}).SetStatusCode(specs.StatusCodeBadGateway)
	}

	req.Hijack(func(conn net.Conn) {
		defer dest_conn.Close()
		
		io.Copy(conn, dest_conn)
		io.Copy(dest_conn, conn)
	})

	return (&giglet.HeaderResponse{}).SetStatusCode(specs.StatusCodeOK)
}
