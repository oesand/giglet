package specs

import (
	"giglet/safe"
	"strings"
)

func parseCookies(cookie string) map[string]string {
	var output map[string]string
	for _, pair := range strings.Split(cookie, "; ") {
		key, value, ok := strings.Cut(pair, "=")
		if !ok || len(key) == 0 || len(value) == 0 {
			continue
		} else if output == nil {
			output = map[string]string{}
		}
		output[key] = value
	}
	return output
}

func NewReadOnlyHeader(headers map[string]string) *ReadOnlyHeader {
	return &ReadOnlyHeader{
		headers: headers,
	}
}

type ReadOnlyHeader struct {
	_ safe.NoCopy

	headers       map[string]string
	cookies       map[string]string
	cookiesParsed bool
}

func (header *ReadOnlyHeader) Get(name string) string {
	if header.headers == nil {
		return ""
	}
	return header.headers[name]
}

func (header *ReadOnlyHeader) GetCookie(name string) string {
	if header.headers == nil {
		return ""
	} else if !header.cookiesParsed {
		if cookie, exists := header.headers["Cookie"]; exists && len(cookie) > 0 {
			header.cookies = parseCookies(cookie)
		}
		header.cookiesParsed = true
	} else if header.cookies == nil {
		return ""
	}
	return header.cookies[name]
}
