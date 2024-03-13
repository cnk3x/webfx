package bufpool

import (
	"bytes"
	"sync"
)

var bufPool = sync.Pool{New: func() any { return new(bytes.Buffer) }}

func Get() *bytes.Buffer {
	b := bufPool.Get().(*bytes.Buffer)
	b.Reset()
	return b
}

func Put(buf *bytes.Buffer) {
	buf.Reset()
	bufPool.Put(buf)
}
