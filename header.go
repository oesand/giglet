package giglet

import (
	"giglet/safe"
	"giglet/specs"
	"io"
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

func (header *Header) Set(name, value string) {
	name = titleCaser.String(name)
	if name == "Set-Cookie" {
		panic("header not support direct set cookie, use method 'SetCookie'")
	} else if header.headers == nil {
		header.headers = map[string]string{}
	}
	header.headers[name] = value
}

func (header *Header) GetCookie(name string) *specs.Cookie {
	if header.cookies == nil {
		return nil
	}
	return header.cookies[name]
}

func (header *Header) SetCookie(cookie *specs.Cookie) {
	if len(cookie.Name) == 0 {
		panic("cookie name cannot be empty")
	} else if header.cookies == nil {
		header.cookies = map[string]*specs.Cookie{}
	}
	
	header.cookies[cookie.Name] = cookie
}

func (header *Header) SetCookieValue(name, value string) {
	header.SetCookie(&specs.Cookie{
		Name: name,
		Value: value,
	})
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
