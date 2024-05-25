package giglet

import "errors"

type RequestHandler func(request *HttpRequest) Response

const (
	defaultLineLengthLimit = 1024
)

var (
	ErrorTooLarge = errors.New("too large")
	ErrorServerClosed = errors.New("http server closed")

	httpVersionPrefix 	= []byte("HTTP/")
	httpV10 			= []byte("HTTP/1.0")
	httpV11 			= []byte("HTTP/1.1")
	httpV2 				= []byte("HTTP/2.0")

	directCrlf              = []byte("\r\n")
	directColonSpace        = []byte(": ")

	headerSetCookie			= []byte("Set-Cookie: ")
)