package ringbuffer

import "encoding/binary"

// READ LOCK
func (this *RingBuffer) Peek(len int) (first []byte, end []byte) {
	this.m.RLock()
	defer this.m.RUnlock()

	if this.isEmpty || len <= 0 {
		return
	}

	if this.wIdx > this.rIdx {
		if len > this.wIdx-this.rIdx {
			len = this.wIdx - this.rIdx
		}

		first = this.buf[this.rIdx : this.rIdx+len]
		return
	}

	if len > this.cap-this.rIdx+this.wIdx {
		len = this.cap - this.rIdx + this.wIdx
	}
	if this.rIdx+len <= this.cap {
		first = this.buf[this.rIdx : this.rIdx+len]
	} else {
		// head
		first = this.buf[this.rIdx:this.cap]
		// tail
		end = this.buf[0 : len-this.cap+this.rIdx]
	}
	return
}

// READ LOCK
func (this *RingBuffer) PeekAll() (first []byte, end []byte) {
	this.m.RLock()
	defer this.m.RUnlock()

	if this.isEmpty {
		return
	}

	if this.wIdx > this.rIdx {
		first = this.buf[this.rIdx:this.wIdx]
		return
	}

	first = this.buf[this.rIdx:this.cap]
	end = this.buf[0:this.wIdx]
	return
}

// READ LOCK
func (this *RingBuffer) PeekUint8() uint8 {
	this.m.RLock()
	defer this.m.RUnlock()

	if this.size() < 1 {
		return 0
	}

	f, e := this.Peek(1)
	if len(e) > 0 {
		return e[0]
	} else {
		return f[0]
	}
}

// READ LOCK
func (this *RingBuffer) PeekUint16() uint16 {
	this.m.RLock()
	defer this.m.RUnlock()

	if this.size() < 2 {
		return 0
	}

	f, e := this.Peek(2)
	if len(e) > 0 {
		return binary.BigEndian.Uint16(bytesJoin2NewByteSlice(f, e))
	} else {
		return binary.BigEndian.Uint16(f)
	}
}

// READ LOCK
func (this *RingBuffer) PeekUint32() uint32 {
	this.m.RLock()
	defer this.m.RUnlock()

	if this.size() < 4 {
		return 0
	}

	f, e := this.Peek(4)
	if len(e) > 0 {
		return binary.BigEndian.Uint32(bytesJoin2NewByteSlice(f, e))
	} else {
		return binary.BigEndian.Uint32(f)
	}
}

// READ LOCK
func (this *RingBuffer) PeekUint64() uint64 {
	this.m.RLock()
	defer this.m.RUnlock()

	if this.size() < 8 {
		return 0
	}

	f, e := this.Peek(8)
	if len(e) > 0 {
		return binary.BigEndian.Uint64(bytesJoin2NewByteSlice(f, e))
	} else {
		return binary.BigEndian.Uint64(f)
	}
}








