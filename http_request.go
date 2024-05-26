package giglet

import (
	"giglet/safe"
	"giglet/specs"
	"giglet/url"
	"io"
	"net"
)

type HttpRequest struct {
	_ safe.NoCopy

	server *Server
	conn net.Conn
	hijacker HijackHandler
	extras map[string]any

	method specs.HttpMethod
	url *url.Url
	header *httpRequestHeader

	protoMajor uint16
	protoMinor uint16

	stream io.ReadCloser
}

func (req *HttpRequest) ProtoAtLeast(major, minor uint16) bool {
	return req.protoMajor > major ||
		req.protoMajor == major && req.protoMinor >= minor
}

func (req *HttpRequest) Stream() io.ReadCloser {
	return req.stream
}

func (req *HttpRequest) RemoteAddr() net.Addr {
	return req.conn.RemoteAddr()
}

func (req *HttpRequest) Hijack(handler HijackHandler) {
	req.hijacker = handler
}

func (req *HttpRequest) Method() specs.HttpMethod {
	return req.method
}

func (req *HttpRequest) Url() *url.Url {
	return req.url
}

func (req *HttpRequest) Header() *httpRequestHeader {
	return req.header
}

func (req *HttpRequest) GetExtra(key string) any {
	if req.extras == nil {
		return nil
	}
	return req.extras[key]
}

func (req *HttpRequest) SetExtra(key string, value any) {
	if req.extras == nil {
		req.extras = map[string]any{}
	}
	req.extras[key] = value
}
