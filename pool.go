package giglet

import (
	"bufio"
	"io"
	"net"
)

func readBufferLine(reader *bufio.Reader, limit int64) ([]byte, error) {
	var line []byte
	for {
		l, more, err := reader.ReadLine()
		if err != nil {
			return nil, err
		} else if limit > 0 && int64(len(line)) + int64(len(l)) > limit {
			return nil, ErrorTooLarge
		} else if line == nil && !more {
			return l, nil
		}
		line = append(line, l...)
		if !more {
			break
		}
	}
	return line, nil
}

func isCommonNetReadError(err error) bool {
	if err == io.EOF {
		return true
	}
	if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
		return true
	}
	if oe, ok := err.(*net.OpError); ok && oe.Op == "read" {
		return true
	}
	return false
}

