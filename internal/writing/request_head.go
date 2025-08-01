package writing

import (
	"bytes"
	"github.com/oesand/giglet/specs"
	"io"
)

func WriteRequestHead(writer io.Writer, method specs.HttpMethod, url *specs.Url, header *specs.Header) (int64, error) {
	// Headline
	buf := bytes.NewBufferString(string(method))
	buf.WriteRune(' ')
	if path := url.Path; path != "" {
		buf.WriteString(path)
	} else {
		buf.WriteByte('/')
	}

	if query := url.Query; query != nil && len(query) > 0 {
		buf.WriteRune('?')
		buf.WriteString(query.String())
	}

	buf.WriteRune(' ')
	buf.Write(httpV11)

	buf.Write(rawCrlf)

	// Headers
	for key, value := range header.All() {
		buf.WriteString(key)
		buf.Write(rawColonSpace)
		buf.WriteString(value)
		buf.Write(rawCrlf)
	}

	if header.AnyCookies() {
		buf.Write(rawCookie)

		firstCookie := true
		for cookie := range header.Cookies() {
			if firstCookie {
				firstCookie = false
			} else {
				buf.Write(rawCookieDelimiter)
			}
			buf.WriteString(cookie.Name)
			buf.WriteRune('=')
			buf.WriteString(cookie.Value)
		}

		buf.Write(rawCrlf)
	}

	buf.Write(rawCrlf)

	return buf.WriteTo(writer)
}
