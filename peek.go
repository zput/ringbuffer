package ringbuffer

import "encoding/binary"

// READ LOCK
func (this *RingBuffer) Peek(len int, isUsingExplore bool) (first []byte, end []byte) {
	this.m.RLock()
	defer this.m.RUnlock()

	var (
		readPosition int
	)

	if isUsingExplore == true {
		if this.episEmpty || len <= 0 {
			return
		}
		readPosition = this.eprIdx
	}else{
		if this.isEmpty || len <= 0 {
			return
		}
		readPosition = this.rIdx
	}

	if this.wIdx > readPosition {
		if len > this.wIdx-readPosition {
			len = this.wIdx - readPosition
		}

		first = this.buf[readPosition : readPosition+len]
		return
	}

	if len > this.cap-readPosition+this.wIdx {
		len = this.cap - readPosition + this.wIdx
	}
	if readPosition+len <= this.cap {
		first = this.buf[readPosition : readPosition+len]
	} else {
		// head
		first = this.buf[readPosition:this.cap]
		// tail
		end = this.buf[0 : len-this.cap+readPosition]
	}
	return
}

// READ LOCK
func (this *RingBuffer) PeekAll(isUsingExplore bool) (first []byte, end []byte) {
	this.m.RLock()
	defer this.m.RUnlock()

	var (
		readPosition int
	)

	if isUsingExplore == true {
		if this.episEmpty {
			return
		}
		readPosition = this.eprIdx
	}else{
		if this.isEmpty {
			return
		}
		readPosition = this.rIdx
	}

	if this.wIdx > readPosition {
		first = this.buf[readPosition:this.wIdx]
		return
	}

	first = this.buf[readPosition:this.cap]
	end = this.buf[0:this.wIdx]
	return
}

// READ LOCK
func (this *RingBuffer) PeekUint8(isUsingExplore bool) uint8 {
	this.m.RLock()
	defer this.m.RUnlock()

	if isUsingExplore == true {
		if this.ExploreSize() < 1 {
			return 0
		}
	}else{
		if this.size() < 1 {
			return 0
		}
	}

	f, e := this.Peek(1, isUsingExplore)
	if len(e) > 0 {
		return e[0]
	} else {
		return f[0]
	}
}

// READ LOCK
func (this *RingBuffer) PeekUint16(isUsingExplore bool) uint16 {
	this.m.RLock()
	defer this.m.RUnlock()

	if isUsingExplore == true {
		if this.ExploreSize() < 2 {
			return 0
		}
	}else{
		if this.size() < 2 {
			return 0
		}
	}

	f, e := this.Peek(2, isUsingExplore)
	if len(e) > 0 {
		return binary.BigEndian.Uint16(bytesJoin2NewByteSlice(f, e))
	} else {
		return binary.BigEndian.Uint16(f)
	}
}

// READ LOCK
func (this *RingBuffer) PeekUint32(isUsingExplore bool) uint32 {
	this.m.RLock()
	defer this.m.RUnlock()

	if isUsingExplore == true {
		if this.ExploreSize() < 4 {
			return 0
		}
	}else{
		if this.size() < 4 {
			return 0
		}
	}

	f, e := this.Peek(4, isUsingExplore)
	if len(e) > 0 {
		return binary.BigEndian.Uint32(bytesJoin2NewByteSlice(f, e))
	} else {
		return binary.BigEndian.Uint32(f)
	}
}

// READ LOCK
func (this *RingBuffer) PeekUint64(isUsingExplore bool) uint64 {
	this.m.RLock()
	defer this.m.RUnlock()

	if isUsingExplore == true {
		if this.ExploreSize() < 8 {
			return 0
		}
	}else{
		if this.size() < 8 {
			return 0
		}
	}

	f, e := this.Peek(8, isUsingExplore)
	if len(e) > 0 {
		return binary.BigEndian.Uint64(bytesJoin2NewByteSlice(f, e))
	} else {
		return binary.BigEndian.Uint64(f)
	}
}
