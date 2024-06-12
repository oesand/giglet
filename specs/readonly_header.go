package specs

import (
	"giglet/safe"
	"mime"
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

// Creates read-only headers struct from mapped valid cased (Title-Case) headers map
func NewReadOnlyHeader(headers map[string]string) *ReadOnlyHeader {
	header := &ReadOnlyHeader{ headers: headers }
	if media, has := headers["Content-Type"]; has {
		contentType, mediaParams, err := mime.ParseMediaType(media)
		if err != nil {
			header.contentType = ContentType(media)
		} else {
			header.contentType = ContentType(contentType)
			header.mediaParams = mediaParams
		}
		delete(headers, "Content-Type")
	}
	return header
}

type ReadOnlyHeader struct {
	_ safe.NoCopy

	contentType 		ContentType
	mediaParams			map[string]string

	headers       		map[string]string
	cookies       		map[string]string
	cookiesParsed 		bool
}

func (header *ReadOnlyHeader) Get(name string) string {
	if header.headers == nil {
		return ""
	} else if name == "Content-Type" {
		return string(header.contentType)
	}

	return header.headers[name]
}

func (header *ReadOnlyHeader) ContentType() ContentType {
	return header.contentType
}

func (header *ReadOnlyHeader) GetMediaParams(name string) string {
	if header.mediaParams == nil {
		return ""
	}
	return header.mediaParams[name]
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
