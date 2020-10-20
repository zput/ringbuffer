#  Multi-function Ringbuffer

[![LICENSE](https://img.shields.io/badge/LICENSE-MIT-blue)](https://github.com/zput/ringbuffer/blob/master/LICENSE)
[![Github Actions](https://github.com/zput/ringbuffer/workflows/CI/badge.svg)](https://github.com/zput/ringbuffer/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/zput/ringbuffer)](https://goreportcard.com/report/github.com/zput/ringbuffer)
[![GoDoc](https://godoc.org/github.com/zput/ringbuffer?status.svg)](https://godoc.org/github.com/zput/ringbuffer)


#### [中文](README-ZH.md) | English

- Control whether locking(thread safe) or unlocking(single thread; fast) is required via parameters
- Automatic expansion of the circular buffer implementation
- Pre-read the data in the cache by exploring(ExploreBegin()----ExploreRead()/ExploreSize()----ExploreCommit()/ExploreBreak())


## Features

- Freedom to decide whether to lock or not, balancing performance and thread safety
- Automatically expands space when cache is full
- Provides peek at cached content in advance
- Provide explore class functions that simulate reading first, but don't move the actual


## Performance Testing

<details>
  <summary> 📈 测试数据 </summary>

> os platform: Mac 

### test for write and read

have locked

```golang
goos: darwin
goarch: amd64
pkg: github.com/zput/ringbuffer
BenchmarkRingBuffer_Sync_Unlock-4   	29223921	        43.5 ns/op
PASS
```

unlocked

```golang
goos: darwin
goarch: amd64
pkg: github.com/zput/ringbuffer
BenchmarkRingBuffer_Sync_Lock-4   	12641550	        89.1 ns/op
PASS
```

</details>


## Example

<details>
  <summary> Create locked/unlocked ringbuffer objects</summary>

```go
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

```

</details>

<details>
  <summary> Peeking into the unread data in this Ringbuffer </summary>

```go
package main

import (
	"fmt"
	"github.com/zput/ringbuffer"
)

const bufferCapacity = 1024

func main() {
	// default not thread safe
	buffer := ringbuffer.New(bufferCapacity)

	fmt.Println(buffer.WriteString("writing ..."))

	fmt.Printf("size[%d]; capacity[%d]\n", buffer.Size(), buffer.Capacity())

	print := func(first, second []byte) {
		if len(second) == 0 {
			fmt.Println(string(first))
		} else {
			first = append(first, second...)
			fmt.Println(string(first))
		}
	}

	print(buffer.Peek(7))

	print(buffer.PeekAll())

	fmt.Println(buffer.PeekUint8())

	fmt.Println(buffer.PeekUint16())

	fmt.Println(buffer.PeekUint32())

	fmt.Println(buffer.PeekUint64())
}

/*
11 <nil>
size[11]; capacity[1024]
writing
writing ...
119
30578
2003986804
8607057786564405024

Process finished with exit code 0
*/
```

</details>

<details>
  <summary> Explore Class Functions </summary>

```go
package main

import (
	"fmt"
	"github.com/zput/ringbuffer"
)

const bufferCapacity = 1024

func main() {
	// default not thread safe
	buffer := ringbuffer.New(bufferCapacity)

	fmt.Println(buffer.WriteString("writing ..."))

	fmt.Println(buffer.PrintRingBufferInfo())

	buffer.ExploreBegin()

	buf := make([]byte, 4)

	n, err := buffer.ExploreRead(buf)
	if err != nil {
		panic(err)
	}

	fmt.Printf("read %d byte through ExploreRead;\nremaining %d size to explore read;\nremaining %d size to read;\n", n, buffer.ExploreSize(), buffer.Size())

	buffer.ExploreCommit()

	fmt.Println("====after commit=====")

	fmt.Printf("remaining %d size to explore read;\nremaining %d size to read;", buffer.ExploreSize(), buffer.Size())

}

/*
11 <nil>

	Ring Buffer:
		Cap: 1024
		size(can read): 11
		FreeSpace: 1013
		Content: writing ...

read 4 byte through ExploreRead;
remaining 7 size to explore read;
remaining 11 size to read;
====after commit=====
remaining 7 size to explore read;
remaining 7 size to read;
Process finished with exit code 0
*/
```

</details>


## Reference

- https://github.com/smallnest/ringbuffer
- https://github.com/Allenxuxu/ringbuffer


## Appendix

welcome pr
