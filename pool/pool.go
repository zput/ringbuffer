package pool

import (
	"sync"

	"github.com/zput/ringbuffer"
)

var DefaultPool = New(1024)

func Get() *ringbuffer.RingBuffer {
	return DefaultPool.Get()
}

func Put(rIdx *ringbuffer.RingBuffer) {
	DefaultPool.Put(rIdx)
}

type RingBufferPool struct {
	pool *sync.Pool
}

func New(initSize int) *RingBufferPool {
	return &RingBufferPool{
		pool: &sync.Pool{
			New: func() interface{} {
				return ringbuffer.New(initSize)
			},
		},
	}
}

/*
 - 从pool中任意返回一个对象，且从池子中删除它；
 - 即使pool中有数据，但是视它为空
 - 如果pool.New字段不是空，且从上面没有得到值，那么会调用New函数新生成一个对象返回。
*/
func (p *RingBufferPool) Get() *ringbuffer.RingBuffer {
	rIdx, _ := p.pool.Get().(*ringbuffer.RingBuffer)
	return rIdx
}

func (p *RingBufferPool) Put(rIdx *ringbuffer.RingBuffer) {
	p.pool.Put(rIdx)
}
