package giglet

import (
	"errors"
	"net"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type RequestHandler func(request *HttpRequest) Response
type HijackHandler func(c net.Conn)
type NextProtoHandler func(c net.Conn) error


var (
	HeadlineMaxLength uint64 = 1024
	HeaderMaxLength uint64 = 1024

	ErrorTooLarge = errors.New("too large")
	ErrorServerClosed = errors.New("http: server closed")
	ErrorUnsupportedEncoding = errors.New("http: encoding not supported")
	ErrorAbortHandler = errors.New("http: abort Handler")

	titleCaser = cases.Title(language.English)
	zeroTime time.Time
	
	httpVersionPrefix 	= []byte("HTTP/")
	httpV10 			= []byte("HTTP/1.0")
	httpV11 			= []byte("HTTP/1.1")
	httpV2 				= []byte("HTTP/2.0")
	
	headerSetCookie			= []byte("Set-Cookie: ")
	directCrlf              = []byte("\r\n")
	directColonSpace        = []byte(": ")
	directColon        		= []byte(":")
	
	responseDowngradeHTTPS 	= []byte("HTTP/1.0 400 Bad Request\r\n\r\nSent an HTTP request to an HTTPS server.\n")
	responseRequestHeadersTooLarge 	= []byte("HTTP/1.1 431 Request Header Fields Too Large\r\nContent-Type: text/plain; charset=utf-8\r\nConnection: close\r\n\r\n431 Request Header Fields Too Large\n")
	responseUnsupportedEncoding 	= []byte("HTTP/1.1 431 Request Header Fields Too Large\r\nContent-Type: text/plain; charset=utf-8\r\nConnection: close\r\n\r\n431 Request Header Fields Too Large\n")
)