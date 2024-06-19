package giglet

import (
	"bytes"
	"crypto/tls"
	"errors"
	"giglet/specs"
	"io"
	"net"
	"time"
)

type RequestHandler func(request Request) Response
type HijackHandler func(conn net.Conn)
type NextProtoHandler func(conn *tls.Conn)
type EventHandler func()

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
	httpV1NextProtoTLS 	= "http/1.1"

	httpVersionPrefix 	= []byte("HTTP/")
	httpV10 			= []byte("HTTP/1.0")
	httpV11 			= []byte("HTTP/1.1")
	httpV2 				= []byte("HTTP/2.0")
	
	directCrlf              = []byte("\r\n")
	directColon        		= []byte(": ")

	rawCloseHeaders 				= []byte("Content-Type: text/plain; charset=utf-8\r\nConnection: close\r\n")
	responseDowngradeHTTPS 			= []byte("HTTP/1.0 400 Bad Request\r\n\r\nSent an HTTP request to an HTTPS server.\n")
	responseRequestHeadersTooLarge 	= []byte("HTTP/1.1 431 Request Header Fields Too Large\r\nContent-Type: text/plain; charset=utf-8\r\nConnection: close\r\n\r\n431 Request Header Fields Too Large\n")
	responseNotProcessableError 	= []byte("HTTP/1.1 500 Internal Server Error\r\nContent-Type: text/plain; charset=utf-8\r\nConnection: close\r\n\r\n500 Unknown error while processing the request\n")
	responseUnsupportedEncoding 	= []byte("HTTP/1.1 501 Not Implemented\r\nContent-Type: text/plain; charset=utf-8\r\nConnection: close\r\n\r\n501 Unsupported transfer encoding\n")
)

type statusErrorResponse struct {
	code specs.StatusCode
	text string
}

func (err *statusErrorResponse) Error() string {
	return string(err.code.AppendBytes(nil)) + ": " + err.text
}

func (err *statusErrorResponse) Write(writer io.Writer) {
	var buf bytes.Buffer

	buf.Write(httpV11)
	buf.WriteByte(' ')
	buf.Write(err.code.AppendBytes(nil))
	buf.Write(directCrlf)
	buf.Write(rawCloseHeaders)
	buf.Write(directCrlf)
	buf.WriteString(err.text)

	buf.WriteTo(writer)
}
