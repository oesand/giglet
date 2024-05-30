package giglet

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"giglet/safe"
	"giglet/specs"
	"strconv"

	"golang.org/x/net/http/httpguts"
)

func readRequest(reader *bufio.Reader) (*HttpRequest, error) {
	line, err := readBufferLine(reader, HeadlineMaxLength);
	if err != nil {
		return nil, err
	}

	method, rawurl, proto, ok := parseHeadline(line)
	if !ok {
		return nil, &statusErrorResponse{ 
			code: specs.StatusCodeRequestURITooLong,
			text: "http: invalid headline",
		}
	}
	
	var protoMajor, protoMinor uint16
	if protoMajor, protoMinor, ok = parseHTTPVersion(proto); !ok {
		return nil, errors.New("http: invalid http version format")
	} else if protoMajor != 1 && (
		protoMajor != 2 || protoMinor != 0 || method != "PRI") {
		return nil, &statusErrorResponse{
			code: specs.StatusCodeNotImplemented,
			text: fmt.Sprintf("http: unsupported http version %d.%d", protoMajor, protoMinor),
		}
	}
	var url *specs.Url
	if url, err = specs.ParseUrl(safe.BufferToString(rawurl)); err != nil {
		return nil, &statusErrorResponse{
			code: specs.StatusCodeMisdirectedRequest,
			text: fmt.Sprintf("http: invalid request url \"%s\"", rawurl),
		}
	}

	req := new(HttpRequest)
	req.method = specs.HttpMethod(method)
	req.protoMajor, req.protoMinor = protoMajor, protoMinor
	req.url = url
	
	header, err := parseHeader(reader)
	if err != nil {
		return nil, err
	}
	req.header = specs.NewReadOnlyHeader(header)

	if req.ProtoAtLeast(1, 1) { // [FIXME]: Add chunked transfer
		if raw := req.header.Get("Transfer-Encoding"); len(raw) > 0 { // !ascii.EqualFold(raw, "chunked")
			return nil, ErrorUnsupportedEncoding
		}
	}

	// RFC 7230, section 5.3: Must treat
	//	GET /index.html HTTP/1.1
	//	Host: www.google.com
	// and
	//	GET http://www.google.com/index.html HTTP/1.1
	//	Host: doesnt matter
	// the same. In the second case, any Host line is ignored.
	if host, has := header["Host"]; 
		has && len(host) > 0 && httpguts.ValidHostHeader(host) {
		req.url.Host = host
	}

	// RFC 7234, section 5.4: Should treat
	if pragma, has := header["Pragma"]; has && pragma == "no-cache" {
		header["Cache-Control"] = "no-cache"
	}
	
	req.stream = reader
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
		line, err := readBufferLine(reader, 0)
		if err != nil {
			return headers, errors.New("header: " + err.Error())
		} else if len(line) == 0 {
			return headers, nil
		} else if key == nil {
			if len(line) < 2 {
				return headers, ErrorHeaderInvalidFormat
			}

			if bytes.HasSuffix(line, directColon) {
				key = line[:len(line) - 1]
			} else {
				key, value, ok := bytes.Cut(line, directColon)
				if !ok || len(key) == 0 || len(value) == 0 {
					continue
				}
				skey, sval := safe.BufferToString(specs.TitleCaseBytes(key)), safe.BufferToString(value)
				if httpguts.ValidHeaderFieldName(skey) && httpguts.ValidHeaderFieldValue(sval) {
					headers[skey] = sval
				}
				headers[skey] = sval
			}
		} else {
			line = bytes.TrimLeft(line, " \t")
			if len(key) == 0 || len(line) == 0 {
				continue
			}
			skey, sval := safe.BufferToString(specs.TitleCaseBytes(key)), safe.BufferToString(line)
			if httpguts.ValidHeaderFieldName(skey) && httpguts.ValidHeaderFieldValue(sval) {
				headers[skey] = sval
			}
			key = nil
		}
	}
}

// parse first line: GET /index.html HTTP/1.0
func parseHeadline(line []byte) (string, []byte, []byte, bool) {
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
		return "", nil, nil, false
	}
	return safe.BufferToString(method), uri, proto, true
}

func parseHTTPVersion(vers []byte) (major, minor uint16, ok bool) {
	if bytes.EqualFold(vers, httpV10) {
		return 1, 0, true
	} else if bytes.EqualFold(vers, httpV11) {
		return 1, 1, true
	} else if bytes.EqualFold(vers, httpV2) {
		return 2, 0, true
	} else if !bytes.HasPrefix(vers, httpVersionPrefix) || 
		len(vers) != 8 || vers[6] != '.' {
		return 0, 0, false
	}

	maj, err := strconv.ParseUint(string(vers[5:6]), 10, 16)
	if err != nil {
		return 0, 0, false
	}
	min, err := strconv.ParseUint(string(vers[7:8]), 10, 16)
	if err != nil {
		return 0, 0, false
	}
	return uint16(maj), uint16(min), true
}
