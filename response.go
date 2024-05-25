package giglet

import (
	"giglet/specs"
	"io"
)

type Response interface {
	StatusCode() *specs.StatusCode
	Header() *Header
	Write(buf io.Writer)
}
