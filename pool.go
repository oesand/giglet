package giglet

import (
	"bufio"
	"io"
	"sync"
)

var (
	readerPool  sync.Pool
	// writerPool 	sync.Pool
)

func getBufloReader(r io.Reader) *bufio.Reader {
	if v := readerPool.Get(); v != nil {
		br := v.(*bufio.Reader)
		br.Reset(r)
		return br
	}
	return bufio.NewReader(r)
}

// func GetBufioSizedWriter(writer io.Writer, size int) *bufio.Writer {
// 	pool := bufioWriterPool(size)
// 	if pool != nil {
// 		if v := pool.Get(); v != nil {
// 			bw := v.(*bufio.Writer)
// 			bw.Reset(w)
// 			return bw
// 		}
// 	}
// 	return bufio.NewWriterSize(w, size)
// }
