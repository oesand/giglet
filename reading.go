package giglet

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"giglet/safe"
	"giglet/specs"
	urlpkg "giglet/url"
)

func readRequest(reader *bufio.Reader) (*HttpRequest, error) {
	line, err := readBufferLine(reader, HeadlineMaxLength);
	if err != nil {
		return nil, err
	}

	method, url, proto, ok := parseHeadline(line)
	if !ok {
		return nil, errors.New("http: invalid headline")
	}

	req := new(HttpRequest)
	req.method = specs.HttpMethod(safe.BufferToString(method))
	if req.protoMajor, req.protoMinor, ok = parseHTTPVersion(proto); !ok {
		return nil, errors.New("http: invalid http version")
	} else if req.url, err = urlpkg.ParseUrl(safe.BufferToString(url)); err != nil {
		return nil, errors.New("parse -> url -> " + err.Error())
	}

	header, err := parseHeader(reader)
	if err != nil {
		return nil, err
	}
	req.header = &httpRequestHeader{
		headers: header,
	}

	// RFC 7230, section 5.3: Must treat
	//	GET /index.html HTTP/1.1
	//	Host: www.google.com
	// and
	//	GET http://www.google.com/index.html HTTP/1.1
	//	Host: doesntmatter
	// the same. In the second case, any Host line is ignored.
	if host, has := header["Host"]; 
		has && len(host) > 0 && len(req.url.Host) == 0 {
		req.url.Host = host
	}

	// RFC 7234, section 5.4: Should treat
	if pragma, has := header["Pragma"]; has && pragma == "no-cache" {
		header["Cache-Control"] = "no-cache"
	}

	return req, nil
}

func parseHeader(reader *bufio.Reader) (map[string]string, error) {
	// The first line cannot start with a leading space.
	if buf, err := reader.Peek(1); err == nil && (buf[0] == ' ' || buf[0] == '\t') {
		line, err := readBufferLine(reader, 50)
		if err != nil {
			return nil, err
		}
		return nil, errors.New("malformed header initial line: " + string(line))
	}

	headers := map[string]string{}

	var key []byte
	for {
		line, err := readBufferLine(reader, HeaderMaxLength)
		if err != nil {
			return headers, errors.New("header: " + err.Error())
		} else if line == nil || len(line) == 0 {
			return headers, nil
		} else if key == nil {
			if len(line) < 2 {
				return headers, errors.New("header: invalid format")
			}

			if line[len(line) - 1] == ':' {
				key = line[:len(line) - 1]
			} else {
				key, value, ok := bytes.Cut(line, directColon)
				if !ok || len(key) == 0 || len(value) == 0 {
					continue
				}
				headers[safe.BufferToString(titleCaser.Bytes(key))] = safe.BufferToString(value)
			}
		} else {
			line = bytes.TrimLeft(line, " \t")
			if len(key) == 0 || len(line) == 0 {
				continue
			}
			headers[safe.BufferToString(titleCaser.Bytes(key))] = safe.BufferToString(line)
			key = nil
		}
	}
	return headers, nil
}

// parse first line: GET /index.html HTTP/1.0
func parseHeadline(line []byte) ([]byte, []byte, []byte, bool) {
	var method, uri, proto []byte
	for i, b := range line {
		if b == ' ' {
			if method == nil {
				method = line[:i]
				if i < 3 { break }
			} else {
				if i - len(method) <= 1 ||
					len(line) - i <= 5 { break }

				uri = line[len(method) + 1:i]
				proto = line[i+1:]
			}
		}
	}
	if method == nil || uri == nil || proto == nil {
		return nil, nil, nil, false
	}
	return method, uri, proto, true
}

func parseHTTPVersion(vers []byte) (major, minor uint16, ok bool) {
	if bytes.Compare(vers, httpV10) == 0 {
		return 1, 0, true
	} else if bytes.Compare(vers, httpV11) == 0 {
		return 1, 1, true
	} else if bytes.Compare(vers, httpV2) == 0 {
		return 2, 0, true
	} else if !bytes.HasPrefix(vers, httpVersionPrefix) || 
		len(vers) != 8 || vers[6] != '.' {
		return 0, 0, false
	}

	major = binary.BigEndian.Uint16(vers[5:6])
	minor = binary.BigEndian.Uint16(vers[7:8])
	ok = major != 0
	return
}
