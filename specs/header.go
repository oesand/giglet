package specs

import (
	"bytes"
	"giglet/safe"
	"io"
)

type Header struct {
	_ safe.NoCopy
	
	headers map[string]string
	cookies map[string]*Cookie
}

func (header *Header) Get(name string) string {
	if header.headers == nil {
		return ""
	}
	return header.headers[name]
}

func (header *Header) Has(name string) bool {
	if header.headers == nil {
		return false
	}
	_, has := header.headers[name]
	return has
}

func (header *Header) Set(name, value string) {
	name = TitleCase(name)
	if name == "Set-Cookie" {
		panic("header not support direct set cookie, use method 'SetCookie'")
	} else if header.headers == nil {
		header.headers = map[string]string{}
	}
	header.headers[name] = value
}

func (header *Header) Del(name string) {
	if header.headers != nil {
		delete(header.headers, TitleCase(name))
	}
}

func (header *Header) GetCookie(name string) *Cookie {
	if header.cookies == nil {
		return nil
	}
	return header.cookies[name]
}

func (header *Header) HasCookie(name string) bool {
	if header.cookies == nil {
		return false
	}
	_, has := header.cookies[name]
	return has
}

func (header *Header) DelCookie(name string) {
	if header.cookies != nil {
		delete(header.cookies, TitleCase(name))
	}
}

func (header *Header) SetCookie(cookie *Cookie) {
	if len(cookie.Name) == 0 {
		panic("cookie name cannot be empty")
	} else if header.cookies == nil {
		header.cookies = map[string]*Cookie{}
	}
	
	header.cookies[cookie.Name] = cookie
}

func (header *Header) SetCookieValue(name, value string) {
	header.SetCookie(&Cookie{
		Name: name,
		Value: value,
	})
}

func (header *Header) WriteTo(writer io.Writer) (int64, error) {
	var buf bytes.Buffer
	
	if header.headers != nil {
		for key, value := range header.headers {
			buf.Write(safe.StringToBuffer(key))
			buf.Write(directColonSpace)
			buf.Write(safe.StringToBuffer(value))
			buf.Write(directCrlf)
		}
	}
	if header.cookies != nil {
		for _, cookie := range header.cookies {
			buf.Write(headerSetCookie)
			buf.Write(cookie.Bytes())
			buf.Write(directCrlf)
		}
	}
	return buf.WriteTo(writer)
}
