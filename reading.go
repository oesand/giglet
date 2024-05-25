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
	line, err := readBufferLine(reader);
	if err != nil {
		return nil, err
	}

	req := new(HttpRequest)

	method, url, proto, ok := parseHeadline(line)
	if !ok {
		return nil, errors.New("parse: invalid headline")
	}

	req.method = specs.HttpMethod(safe.BufferToString(method))
	if !req.method.IsValid() {
		return nil, errors.New("parse: invalid http method")
	} else if req.protoMajor, req.protoMinor, ok = parseHTTPVersion(proto); !ok {
		return nil, errors.New("parse: invalid http version")
	} else if req.url, err = urlpkg.ParseUrl(safe.BufferToString(url)); err != nil {
		return nil, errors.New("parse -> url -> " + err.Error())
	}

	
	



	return req, nil
}

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
