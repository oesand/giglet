package giglet

import (
	"bufio"
	"giglet/safe"
	"net"
)

type HttpRequest struct {
	nocopy safe.NoCopy

	conn net.Conn
	reader bufio.Reader

	
}

func (req *HttpRequest) readLine() ([]byte, error) {
	var line []byte
	for {
		l, more, err := req.reader.ReadLine()
		if err != nil {
			return nil, err
		} else if len(line) + len(l) > defaultLineLengthLimit {
			return nil, errorTooLarge
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