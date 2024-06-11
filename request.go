package giglet

import (
	"giglet/specs"
	"io"
	"mime/multipart"
	"net"
)

type Request interface {
	GetExtra(key string) any
	SetExtra(key string, value any)
	
	ProtoAtLeast(major, minor uint16) bool
	RemoteAddr() net.Addr
	Hijack(handler HijackHandler)
	
	Method() specs.HttpMethod
	Url() *specs.Url
	Header() *specs.ReadOnlyHeader
	Stream() io.Reader
	PostForm() (specs.Query, error)
	MultipartForm() (*multipart.Form, error)
}
