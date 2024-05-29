package giglet

import (
	"bytes"
	"errors"
	"giglet/specs"
	"io"
	"net"
	"time"
)

type RequestHandler func(request *HttpRequest) Response
type HijackHandler func(c net.Conn)
type NextProtoHandler func(c net.Conn) error


var (
	DefaultServerName = "giglet"
	HeadlineMaxLength int64 = 2048
	DefaultContentMaxSizeBytes int64 = 5 << 20 // 5MB

	ErrorTooLarge = errors.New("too large")
	ErrorServerClosed = errors.New("http: server closed")
	ErrorHeaderInvalidFormat = errors.New("header: invalid format")
	ErrorUnsupportedEncoding = errors.New("http: encoding not supported")
	ErrorAbortHandler = errors.New("http: abort Handler")

	zeroTime time.Time
	
	httpVersionPrefix 	= []byte("HTTP/")
	httpV10 			= []byte("HTTP/1.0")
	httpV11 			= []byte("HTTP/1.1")
	httpV2 				= []byte("HTTP/2.0")
	
	directCrlf              = []byte("\r\n")
	directColon        		= []byte(":")

	rawCloseHeaders = []byte("Content-Type: text/plain; charset=utf-8\r\nConnection: close\r\n")
	responseDowngradeHTTPS 			= []byte("HTTP/1.0 400 Bad Request\r\n\r\nSent an HTTP request to an HTTPS server.\n")
	responseRequestHeadersTooLarge 	= []byte("HTTP/1.1 431 Request Header Fields Too Large\r\nContent-Type: text/plain; charset=utf-8\r\nConnection: close\r\n\r\n431 Request Header Fields Too Large\n")
	responseNotProcessableError 	= []byte("HTTP/1.1 500 Internal Server Error\r\nContent-Type: text/plain; charset=utf-8\r\nConnection: close\r\n\r\n500 Unknown error while processing the request\n")
	responseUnsupportedEncoding 	= []byte("HTTP/1.1 501 Not Implemented\r\nContent-Type: text/plain; charset=utf-8\r\nConnection: close\r\n\r\n501 Unsupported transfer encoding\n")
)

type statusErrorResponse struct {
	code *specs.StatusCode
	text string
}

func (err *statusErrorResponse) Error() string {
	return string(err.code.Append(nil)) + ": " + err.text
}

func (err *statusErrorResponse) Write(writer io.Writer) {
	buff := bytes.Buffer{}

	buff.Write(httpV11)
	buff.WriteByte(' ')
	buff.Write(err.code.Append(nil))
	buff.Write(directCrlf)
	buff.Write(rawCloseHeaders)
	buff.Write(directCrlf)
	buff.WriteString(err.text)

	buff.WriteTo(writer)
}
