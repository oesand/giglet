package giglet

import (
	"github.com/oesand/giglet/specs"
	"io"
)

type Response interface {
	StatusCode() specs.StatusCode
	SetStatusCode(specs.StatusCode)
	Header() *specs.Header
}

type BodyWriter interface {
	WriteBody(io.Writer) error
}
