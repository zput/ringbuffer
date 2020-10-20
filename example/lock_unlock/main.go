package main

import (
	"fmt"
	"github.com/zput/ringbuffer"
)

const bufferCapacity = 5

func main() {
	var (
		unLockBuffer *ringbuffer.RingBuffer
		lockBuffer *ringbuffer.RingBuffer
	)

	// default not thread safe
	unLockBuffer = ringbuffer.New(bufferCapacity)
	fmt.Println(unLockBuffer.WriteString("writing ..."))
	fmt.Println(unLockBuffer.Size(), unLockBuffer.Capacity())
	fmt.Println(string(unLockBuffer.ReadAll2NewByteSlice()))

	var(
		whetherThreadSafe = true
		data = make([]byte, bufferCapacity, bufferCapacity+2)
		err error
	)

	lockBuffer, err = ringbuffer.NewWithDataAndPointer(data, 0, 0, false, whetherThreadSafe)
	if err != nil{
		panic(err)
	}
	// should equal true
	fmt.Println(lockBuffer.IsFull())
	// size == 5  capacity == 5
	fmt.Println(lockBuffer.Size(), lockBuffer.Capacity())
	// [0 0 0 0 0]
	fmt.Println(lockBuffer.ReadAll2NewByteSlice())

	err = lockBuffer.WriteOneByte(byte(15))
	if err != nil{
		panic(err)
	}
	// [15 0 0 0 0 0] -compare- [15 0 0 0 0]
	// still use same memory between data and lockBuffer
	fmt.Println(lockBuffer.ReadAll2NewByteSlice(), "-compare-", data)
}

/*
11 <nil>
11 11
writing ...
true
5 5
[0 0 0 0 0]
[15 0 0 0 0 0] -compare- [15 0 0 0 0]

Process finished with exit code 0
*/
