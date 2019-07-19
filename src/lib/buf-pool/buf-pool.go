package bufPool

import (
	"bytes"
	"sync"
)

var pool *sync.Pool = nil

type Buf struct {
	bytes.Buffer
	pool *sync.Pool
}

func (it *Buf) Release() {
	if it == nil {
		return
	}
	it.pool.Put(it)
}

func newBuf(pool *sync.Pool) *Buf {
	it := new(Buf)
	it.pool = pool
	return it
}

func NewBuf() *Buf {
	it := pool.Get().(*Buf)
	it.Reset()
	return it
}

func init() {
	pool = &sync.Pool{
		New: func() interface{} {
			buf := newBuf(pool)
			return buf
		},
	}
}
