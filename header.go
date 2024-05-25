package giglet

import (
	"giglet/safe"
	"giglet/specs"
	"io"
	"strings"
)

type Header struct {
	headers map[string]string
	cookies map[string]*specs.Cookie
}

func (header *Header) Get(name string) string {
	if header.headers == nil {
		return ""
	}
	return header.headers[name]
}

func (header *Header) Set(name, value string) *Header {
	name = strings.Title(name)
	if name == "Set-Cookie" {
		panic("header not support direct set cookie, use method 'SetCookie'")
	} else if header.headers == nil {
		header.headers = map[string]string{}
	}
	header.headers[name] = value
	return header
}

func (header *Header) SetCookie(cookie *specs.Cookie) *Header {
	if len(cookie.Name) == 0 {
		panic("cookie name cannot be empty")
	} else if header.cookies == nil {
		header.cookies = map[string]*specs.Cookie{}
	}
	
	header.cookies[cookie.Name] = cookie
	return header
}

func (header *Header) SetCookieValue(name, value string) *Header {
	header.SetCookie(&specs.Cookie{
		Name: name,
		Value: value,
	})
	return header
}

func (header *Header) Write(writer io.Writer) {
	if header.headers != nil {
		for key, value := range header.headers {
			writer.Write(safe.StringToBuffer(key))
			writer.Write(directColonSpace)
			writer.Write(safe.StringToBuffer(value))
			writer.Write(directCrlf)
		}
	}
	if header.cookies != nil {
		for _, cookie := range header.cookies {
			writer.Write(cookie.Append(headerSetCookie))
			writer.Write(directCrlf)
		}
	}
}
