package ringbuffer

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"unsafe"
)

const GrowthFactor float32 = 1.5

type innerLock struct {
	sync.RWMutex
	IsOpen bool
}

func (this *innerLock) RLock() {
	if this.IsOpen {
		this.RWMutex.RLock()
	}
}
func (this *innerLock) RUnlock() {
	if this.IsOpen {
		this.RWMutex.RUnlock()
	}
}
func (this *innerLock) Lock() {
	if this.IsOpen {
		this.RWMutex.Lock()
	}
}
func (this *innerLock) Unlock() {
	if this.IsOpen {
		this.RWMutex.Unlock()
	}
}

// 缓冲区中没有数据：ErrIsEmpty
var ErrIsEmpty = errors.New("ring buffer is empty")

var ErrIsNotInExplore = errors.New("not begin explore read; ring buffer")
var ErrExploreRetrievingCrossTheLine = errors.New("retrieving cross the line in explore mode")

var ErrInitRingBufferParameter = errors.New("parameter is not right; when initializing ring buffer")

/*
	 _ _ _ _ _
	|_|_|_|_|_|
	default: rIdx == wIdx == 0
	 _ _ _ _ _
	|x|_|_|_|_|
	rIdx == 0; array[wIdx]=x; ++wIdx == 1
	 _ _ _ _ _
	|x|_|_|_|_|
	rIdx == 0; wIdx == 1

	in here: write index in array is not saved data, it is free space.
	RingBuffer 循环缓冲区:
      - 当缓冲区满了以后，自动申请新的slice;
	  - 然后把老的slice copy过去。
*/
type RingBuffer struct {
	buf       []byte
	cap       int //这个是缓存的容量; 缓存现保存的数据大小是通过wIdx与rIdx算出来的。
	eprIdx    int
	episEmpty bool
	inExplore bool
	rIdx      int // next position to read
	wIdx      int // next position to write
	isEmpty   bool

	m innerLock
}

// New 返回一个初始大小为 cap 的 RingBuffer
func New(cap int, isOpenLock ...bool) *RingBuffer {
	var isOpen bool
	if len(isOpenLock) > 0 {
		isOpen = isOpenLock[0]
	}
	return &RingBuffer{
		buf:       make([]byte, cap),
		cap:       cap,
		isEmpty:   true,
		m:         innerLock{IsOpen: isOpen},
	}
}

// NewWithData;let ringBuffer is full,
/*
	 _ _ _ _ _
	|0|1|2|3|4|
	rIdx == 0; wIdx == 0; isEmpty == false; this ringBuffer is full
*/
func NewWithData(data []byte, isOpenLock ...bool) *RingBuffer {
	var isOpen bool
	if len(isOpenLock) > 0 {
		isOpen = isOpenLock[0]
	}
	return &RingBuffer{
		buf: data,
		cap: len(data),
		isEmpty:   false,
		m:   innerLock{IsOpen: isOpen},
	}
}

func NewWithDataAndPointer(data []byte, beginPointer, endPointer int, isEmpty bool, isOpenLock ...bool)(*RingBuffer, error) {

	if beginPointer < endPointer && isEmpty != false{
		if len(data)<endPointer{
			return nil, ErrInitRingBufferParameter
		}
		//isEmpty = false
	}else if beginPointer > endPointer && isEmpty != false{
		if len(data)<beginPointer{
			return nil, ErrInitRingBufferParameter
		}
		//isEmpty = false
	}else{
		// beginPointer is equality endPointer
	}

	var isOpen bool
	if len(isOpenLock) > 0 {
		isOpen = isOpenLock[0]
	}
	return &RingBuffer{
		buf: data,
		cap: len(data),
		rIdx:beginPointer,
		wIdx:endPointer,
		isEmpty:isEmpty,
		m:   innerLock{IsOpen: isOpen},
	}, nil
}

// 注意，这个array[wIdx]是没有保存数据的，所以计算剩余空间和已占有空间的时候要注意。
// READ LOCK
// called by inside;  non lock
func (this *RingBuffer) free() int {
	if this.wIdx == this.rIdx {
		if this.isEmpty {
			return this.cap
		}
		return 0
	}

	if this.wIdx < this.rIdx {
		return this.rIdx - this.wIdx
	}

	return this.cap - this.wIdx + this.rIdx
}

// called by inside;  non lock
func (this *RingBuffer) appendSpace(len int) {
	if cap(this.buf) >= this.cap+len{
		reflect.ValueOf(&this.buf).Elem().SetLen(this.cap+len)

		if this.wIdx <= this.rIdx{
			if this.wIdx > (this.cap - this.rIdx) {
				// move from rIdx ->rightIndex
				for i := this.cap - 1; i >= this.rIdx; i-- {
					this.buf[i+len] = this.buf[i]
				}
				this.rIdx += len
			}else{
				// move from 0 -> wIdx
				for i:= 0; i<this.wIdx; i++{
					this.buf[this.cap+i] = this.buf[i]
				}
				this.wIdx = this.cap + this.wIdx
			}
		}
		this.cap += len
	}else{
		newSize := NotMoreThan(this.cap + len)
		//添加增长因子
		//newSize = int(float32(newSize) * GrowthFactor)
		newBuf := make([]byte, newSize)
		oldLen := this.size()
		_, _ = this.read(newBuf)

		this.wIdx = oldLen
		this.rIdx = 0
		this.cap = newSize
		this.buf = newBuf
	}
}

// called by inside;  non lock
func (this *RingBuffer) size() int {
	if this.wIdx == this.rIdx {
		if this.isEmpty {
			return 0
		}
		return this.cap
	}

	if this.wIdx > this.rIdx {
		return this.wIdx - this.rIdx
	}

	return this.cap - this.rIdx + this.wIdx
}

// called by inside;  non lock
func bytesJoin2NewByteSlice(f, e []byte) []byte {
	buf := make([]byte, len(f)+len(e))
	_ = copy(buf, f)
	_ = copy(buf[len(f):], e)
	return buf
}

// called by inside;  non lock
func (this *RingBuffer) read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	if this.isEmpty {
		return 0, ErrIsEmpty
	}
	n = len(p)

	if this.wIdx > this.rIdx {
		/*
		   	 _ _ _ _ _
		   	|x|&|@|_|_|
		        0 1 2 3 4
		   	rIdx==0; wIdx==3
		*/
		if n > this.wIdx-this.rIdx {
			n = this.wIdx - this.rIdx
		}
		copy(p, this.buf[this.rIdx:this.rIdx+n])
		// move readPtr
		this.rIdx = (this.rIdx + n) % this.cap
		if this.rIdx == this.wIdx {
			this.isEmpty = true
		}
		return
	}
	//如果需要读取的数据大于缓存中有的数据，调整n大小等于缓存中的数据长度
	if n > this.cap-this.rIdx+this.wIdx {
		n = this.cap - this.rIdx + this.wIdx
	}
	if this.rIdx+n <= this.cap {
		/*
		   	 _ _ _ _ _
		   	|x|&|@|*|_|
		     0 1 2 3 4
		   	rIdx==0; wIdx==4; if n==1 -->
		*/
		copy(p, this.buf[this.rIdx:this.rIdx+n])
	} else {
		// copy head
		copy(p, this.buf[this.rIdx:this.cap])
		// copy tail
		copy(p[this.cap-this.rIdx:], this.buf[0:n-this.cap+this.rIdx])
	}

	//move read index pointer
	this.rIdx = (this.rIdx + n) % this.cap
	if this.rIdx == this.wIdx {
		this.isEmpty = true
	}
	return
}

func (this *RingBuffer) Capacity() int {
	this.m.RLock()
	defer this.m.RUnlock()

	return this.cap
}

// READ LOCK
func (this *RingBuffer) Size() int {
	this.m.RLock()
	defer this.m.RUnlock()

	return this.size()
}

// READ/WRITE LOCK
func (this *RingBuffer) WriteOneByte(c byte) error {
	this.m.Lock()
	defer this.m.Unlock()

	if this.free() < 1 {
		this.appendSpace(1)
	}

	this.buf[this.wIdx] = c
	this.wIdx++

	if this.wIdx == this.cap {
		this.wIdx = 0
	}

	this.isEmpty = false
	return nil
}

// READ/WRITE LOCK
func (this *RingBuffer) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	this.m.Lock()
	n, err = this.read(p)
	this.m.Unlock()

	return
}

// READ/WRITE LOCK
func (this *RingBuffer) ReadOneByte() (b byte, err error) {
	this.m.Lock()
	defer this.m.Unlock()

	if this.isEmpty {
		return 0, ErrIsEmpty
	}
	b = this.buf[this.rIdx]
	this.rIdx++
	if this.rIdx == this.cap {
		this.rIdx = 0
	}

	if this.wIdx == this.rIdx {
		this.isEmpty = true
	}
	return
}

// READ/WRITE LOCK
func (this *RingBuffer) Write(p []byte) (n int, err error) {

	if len(p) == 0 {
		return 0, nil
	}

	this.m.Lock()
	//defer this.m.Unlock()

	n = len(p)
	free := this.free()
	if free < n {
		this.appendSpace(n - free)
	}
	if this.wIdx >= this.rIdx {
		if this.cap-this.wIdx >= n {
			copy(this.buf[this.wIdx:], p)
			this.wIdx += n
		} else {
			copy(this.buf[this.wIdx:], p[:this.cap-this.wIdx])
			copy(this.buf[0:], p[this.cap-this.wIdx:])
			this.wIdx += n - this.cap
		}
	} else {
		copy(this.buf[this.wIdx:], p)
		this.wIdx += n
	}

	if this.wIdx == this.cap {
		this.wIdx = 0
	}
	this.isEmpty = false

	this.m.Unlock()
	return
}

// non lock; this function calls Write
func (this *RingBuffer) WriteString(s string) (n int, err error) {
	/*	type = struct string {
		    uint8 *str;
		    int len;
		}
		type = struct []uint8 {
		    uint8 *array;
		    int len;
		    int cap;
		}
		string可看做[2]uintptr，而[]byte则是[3]uintptr
	*/
	sPtr := (*[2]uintptr)(unsafe.Pointer(&s))
	u := [3]uintptr{sPtr[0], sPtr[1], sPtr[1]}
	return this.Write(*(*[]byte)(unsafe.Pointer(&u)))
}

/*
   - ReadAll2NewByteSlice:
     - 创建一个新的字节切片，然后把缓存中的所有数据拷贝到新的字节切片去;
     - 不影响原有的缓存任何东西。
*/
// READ LOCK
func (this *RingBuffer) ReadAll2NewByteSlice() (buf []byte) {
	this.m.RLock()
	defer this.m.RUnlock()

	if this.wIdx == this.rIdx {
		if !this.isEmpty {
			buf := make([]byte, this.cap)
			copy(buf, this.buf)
			return buf
		}
		return
	}

	if this.wIdx > this.rIdx {
		buf = make([]byte, this.wIdx-this.rIdx)
		copy(buf, this.buf[this.rIdx:this.wIdx])
		return
	}

	buf = make([]byte, this.cap-this.rIdx+this.wIdx)
	copy(buf, this.buf[this.rIdx:this.cap])
	copy(buf[this.cap-this.rIdx:], this.buf[0:this.wIdx])
	return
}

// READ LOCK
func (this *RingBuffer) IsFull() (isFull bool) {
	this.m.RLock()
	isFull = !this.isEmpty && this.wIdx == this.rIdx
	this.m.RUnlock()
	return
}

// READ LOCK
func (this *RingBuffer) IsEmpty()(isEmpty bool){
	this.m.RLock()
	isEmpty = this.isEmpty
	this.m.RUnlock()
	return
}

// call RetrieveAll
func (this *RingBuffer) Reset() {
	this.RetrieveAll()
}

// READ/WRITE LOCK
func (this *RingBuffer) retrieveAll() {
	this.rIdx = 0
	this.wIdx = 0
	this.isEmpty = true
	this.eprIdx = 0
	this.episEmpty = true
	this.inExplore = false
}

// READ/WRITE LOCK
func (this *RingBuffer) RetrieveAll() {
	this.m.Lock()
	defer this.m.Unlock()

	this.retrieveAll()
}

func (this *RingBuffer) Retrieve(len int) {
	this.m.Lock()

	if this.isEmpty || len <= 0 {
		return
	}

	if len < this.size() {
		this.rIdx = (this.rIdx + len) % this.cap
		if this.wIdx == this.rIdx {
			this.isEmpty = true
		}
	} else {
		this.retrieveAll()
	}
	this.m.Unlock()
}

func (this *RingBuffer) PrintRingBufferInfo() string {
	return fmt.Sprintf("\n\tRing Buffer: \n\t\tCap: %d\n\t\tsize(can read): %d\n\t\tFreeSpace: %d\n\t\tContent: %s\n", this.cap, this.size(), this.free(), this.ReadAll2NewByteSlice())
}

// call ReadOneByte
func (this *RingBuffer) ReadByte() (byte, error) {
	return this.ReadOneByte()
}

// call WriteOneByte
func (this *RingBuffer) WriteByte(c byte) error {
	return this.WriteOneByte(c)
}

/*
	Explore系列的函数，是为了探索缓存中有的数据。

	ExploreBegin
	ExploreRead
	...
	ExploreSize
	ExploreCommit/ExploreBreak
*/
// no thread safety guarantees
func (this *RingBuffer) ExploreBegin() {
	this.eprIdx = this.rIdx
	this.episEmpty = this.isEmpty
	this.inExplore = true
}

func (this *RingBuffer) ExploreCommit() {
	this.rIdx = this.eprIdx
	if this.rIdx == this.wIdx {
		this.isEmpty = true
	}
	this.inExplore = false
}

func (this *RingBuffer) ExploreBreak() {
	this.eprIdx = this.rIdx
	this.episEmpty = this.isEmpty
	this.inExplore = false
}

func (this *RingBuffer) ExploreRetrieve(len int)(err error) {
	if this.inExplore == false {
		return ErrIsNotInExplore
	}

	if this.episEmpty || len <= 0 {
		return
	}

	if len < this.ExploreSize() {
		this.eprIdx = (this.eprIdx + len) % this.cap
		if this.wIdx == this.eprIdx {
			this.episEmpty = true
		}
	} else {
		return ErrExploreRetrievingCrossTheLine
	}
	return nil
}

func (this *RingBuffer) ExploreRead(p []byte) (n int, err error) {
	if this.inExplore == false {
		return 0, ErrIsNotInExplore
	}

	if len(p) == 0 {
		return 0, nil
	}
	if this.episEmpty {
		return 0, ErrIsEmpty
	}
	n = len(p)
	if this.wIdx > this.eprIdx {
		if n > this.wIdx-this.eprIdx {
			n = this.wIdx - this.eprIdx
		}
		copy(p, this.buf[this.eprIdx:this.eprIdx+n])
		// move eprIdx
		this.eprIdx = (this.eprIdx + n) % this.cap
		if this.eprIdx == this.wIdx {
			this.episEmpty = true
		}
		return
	}
	if n > this.cap-this.eprIdx+this.wIdx {
		n = this.cap - this.eprIdx + this.wIdx
	}
	if this.eprIdx+n <= this.cap {
		copy(p, this.buf[this.eprIdx:this.eprIdx+n])
	} else {
		// head
		copy(p, this.buf[this.eprIdx:this.cap])
		// tail
		copy(p[this.cap-this.eprIdx:], this.buf[0:n-this.cap+this.eprIdx])
	}

	// move eprIdx
	this.eprIdx = (this.eprIdx + n) % this.cap
	if this.eprIdx == this.wIdx {
		this.episEmpty = true
	}
	return
}

func (this *RingBuffer) ExploreSize() int {
	if this.wIdx == this.eprIdx {
		if this.episEmpty {
			return 0
		}
		return this.cap
	}

	if this.wIdx > this.eprIdx {
		return this.wIdx - this.eprIdx
	}

	return this.cap - this.eprIdx + this.wIdx
}

const Is64bitArch = ^uint(0) >> 63 == 1
const Is32bitArch = ^uint(0) >> 63 == 0
const WordBits = 32 << (^uint(0) >> 63) // 64或32
const intWordHeadBit = 1 << (WordBits - 2) // 64, 32

// CeilToPowerOfTwo returns the least power of two integer value greater than
// or equal to n.
func NotMoreThan(n int) int {
	if n&intWordHeadBit != 0 && n > intWordHeadBit {
		panic("argument is too large")
	}
	if n <= 2 {
		return 2
	}
	n--
	n = roundingBinaryMath(n)
	n++
	return n
}

/*
0 --- 0
1 --- 1
2 --- 3
3 --- 3
4 --- 7
5 --- 7
6 --- 7
7 --- 7
8 --- 15
*/
func roundingBinaryMath(n int)int{
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n |= n >> 32
	return n
}
