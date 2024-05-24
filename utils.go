package giglet

import "bufio"

func readBufferLine(reader *bufio.Reader) ([]byte, error) {
	var line []byte
	for {
		l, more, err := reader.ReadLine()
		if err != nil {
			return nil, err
		} else if len(line) + len(l) > defaultLineLengthLimit {
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
