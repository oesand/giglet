package writing

import (
	"bytes"
	"github.com/oesand/giglet/specs"
	"io"
	"strconv"
)

func WriteResponseHead(writer io.Writer, is11 bool, code specs.StatusCode, header *specs.Header) (int64, error) {
	if !code.IsValid() {
		code = specs.StatusCodeOK
	}

	// Headline
	var buf bytes.Buffer
	if is11 {
		buf.Write(httpV11)
	} else {
		buf.Write(httpV10)
	}

	buf.WriteRune(' ')
	buf.Write(strconv.AppendUint(nil, uint64(code), 10))
	buf.WriteRune(' ')
	buf.Write(code.Detail())

	buf.Write(rawCrlf)

	// Headers
	for key, value := range header.All() {
		buf.WriteString(key)
		buf.Write(rawColonSpace)
		buf.WriteString(value)
		buf.Write(rawCrlf)
	}

	for cookie := range header.Cookies() {
		buf.Write(rawSetCookie)
		buf.Write(SetCookieBytes(&cookie))
		buf.Write(rawCrlf)
	}

	buf.Write(rawCrlf)

	i, err := buf.WriteTo(writer)
	if err != nil {
		return -1, &specs.GigletError{
			Op:  "write",
			Err: err,
		}
	}
	return i, nil
}
