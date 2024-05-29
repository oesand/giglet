package specs

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const TimeFormat = "Mon, 02 Jan 2006 15:04:05 GMT"

var (
	cookieDelimiter    		= []byte("; ")
	cookieKeyExpires        = []byte("Expires")
	cookieKeyDomain         = []byte("Domain")
	cookieKeyPath           = []byte("Path")
	cookieKeyHTTPOnly       = []byte("HttpOnly")
	cookieKeySecure         = []byte("Secure")
	cookieKeyMaxAge         = []byte("Max-Age")
	cookieKeySameSite       = []byte("SameSite")

	websocketAcceptBaseKey = []byte("258EAFA5-E914-47DA-95CA-C5AB0DC85B11")

	titleCaser 				= cases.Title(language.English)
	directColonSpace        = []byte(": ")
	directCrlf              = []byte("\r\n")
	headerSetCookie			= []byte("Set-Cookie: ")

	httpV10 			= []byte("HTTP/1.0")
	httpV11 			= []byte("HTTP/1.1")
)

func TitleCase(content string) string {
	return titleCaser.String(content)
}

func TitleCaseBytes(content []byte) []byte {
	return titleCaser.Bytes(content)
}
