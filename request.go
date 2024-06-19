package giglet

import (
	"giglet/specs"
	"mime/multipart"
	"net"
)

type Request interface {
	GetExtra(key string) any
	SetExtra(key string, value any)
	
	ProtoAtLeast(major, minor uint16) bool
	ProtoNoHigher(major, minor uint16) bool
	RemoteAddr() net.Addr
	Hijack(handler HijackHandler)
	
	Method() specs.HttpMethod
	Url() *specs.Url
	Header() *specs.ReadOnlyHeader

	Read([]byte) (int, error)
	PostBody() ([]byte, error)
	PostForm() (specs.Query, error)
	MultipartForm() (*multipart.Form, error)
}
