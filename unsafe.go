package giglet

import "unsafe"

func unsafeStringToBuffer(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func unsafeBufferToString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}