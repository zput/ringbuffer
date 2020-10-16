package pool

import (
	"github.com/zput/ringbuffer"
	"testing"
)

func TestRingBufferPool(t *testing.T) {
	pool := New(1024)

	rIdx := pool.Get()
	if rIdx.Capacity() != 1024 {
		t.Fatal()
	}
	if rIdx.Size() != 0 {
		t.Fatal()
	}
	_, _ = rIdx.Write([]byte("1234"))
	pool.Put(rIdx)

	rr := pool.Get()
	if rr.Capacity() != 1024 {
		t.Fatal()
	}
	if rr.Size() != 4 {
		t.Fatal()
	}

	pool.Put(ringbuffer.New(10))
	rrr := pool.Get()
	if rrr.Capacity() != 10 {
		t.Fatal()
	}
	if rrr.Size() != 0 {
		t.Fatal()
	}

	rr.Reset()
	pool.Put(rr)

	rr2 := pool.Get()
	if rr2.Capacity() != 1024 {
		t.Fatal()
	}
	if rr2.Size() != 0 {
		t.Fatal()
	}
}

func TestDefaultPool(t *testing.T) {
	rIdx := Get()
	if rIdx.Capacity() != 1024 {
		t.Fatal()
	}
	if rIdx.Size() != 0 {
		t.Fatal()
	}
	_, _ = rIdx.Write([]byte("1234"))
	Put(rIdx)

	rr := Get()
	if rr.Capacity() != 1024 {
		t.Fatal()
	}
	if rr.Size() != 4 {
		t.Fatal()
	}

	Put(ringbuffer.New(10))
	rrr := Get()
	if rrr.Capacity() != 10 {
		t.Fatal()
	}
	if rrr.Size() != 0 {
		t.Fatal()
	}
}
