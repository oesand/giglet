package giglet

import "net"

type HttpRequest struct {
	nocopy noCopy

	conn net.Conn
}
