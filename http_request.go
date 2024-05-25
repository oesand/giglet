package giglet

import (
	"giglet/safe"
	"giglet/specs"
	"giglet/url"
	"net"
)

type HttpRequest struct {
	_ safe.NoCopy

	server *Server
	conn net.Conn

	method specs.HttpMethod
	url *url.Url

	protoMajor uint16
	protoMinor uint16
}

func (req *HttpRequest) Server() *Server {
	return req.server
}

func (req *HttpRequest) Conn() net.Conn {
	return req.conn
}

func (req *HttpRequest) Method() specs.HttpMethod {
	return req.method
}

func (req *HttpRequest) Url() *url.Url {
	return req.url
}


