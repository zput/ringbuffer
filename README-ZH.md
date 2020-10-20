# 多功能的环形缓存

[![Github Actions](https://github.com/Allenxuxu/gev/workflows/CI/badge.svg)](https://github.com/Allenxuxu/gev/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/Allenxuxu/gev)](https://goreportcard.com/report/github.com/Allenxuxu/gev)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/a2a55fe9c0c443e198f588a6c8026cd0)](https://www.codacy.com/manual/Allenxuxu/gev?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=Allenxuxu/gev&amp;utm_campaign=Badge_Grade)
[![GoDoc](https://godoc.org/github.com/Allenxuxu/gev?status.svg)](https://godoc.org/github.com/Allenxuxu/gev)
[![LICENSE](https://img.shields.io/badge/LICENSE-MIT-blue)](https://github.com/Allenxuxu/gev/blob/master/LICENSE)
[![Code Size](https://img.shields.io/github/languages/code-size/Allenxuxu/gev.svg?style=flat)](https://img.shields.io/github/languages/code-size/Allenxuxu/gev.svg?style=flat)

#### 中文 | [English](README.md)

- 多功能环形缓存：
  - 在New构造函数的时候，通过参数决定是加锁（线程安全）还是不加锁。
  - 当环形缓存空间满后，可以自动扩展内存。
  - 当使用探索类函数(ExploreBegin()----ExploreRead()/ExploreSize()----ExploreCommit()/ExploreBreak())，可以预先探索缓存中的数据，最后可以决定是提交还是放弃。
         
## 特点

- 自由决定是否加锁，在性能和线程安全中取得平衡
- 当缓存空间满，可自动扩展空间
- 提供预先查看缓存中内容（peek）
- 提供探索类函数，可以先模拟读，但是不移动实际的

## 性能测试

<details>
  <summary> 📈 测试数据 </summary>

> 测试电脑 Mac 

### 读写测试

无锁

```golang
goos: darwin
goarch: amd64
pkg: github.com/zput/ringbuffer
BenchmarkRingBuffer_Sync_Unlock-4   	29223921	        43.5 ns/op
PASS
```

有锁

```golang
goos: darwin
goarch: amd64
pkg: github.com/zput/ringbuffer
BenchmarkRingBuffer_Sync_Lock-4   	12641550	        89.1 ns/op
PASS
```

</details>

## 安装

```bash
go get -u github.com/zput/ringbuffer
```

## 示例

<details>
  <summary> 创建有锁/无锁ringbuffer对象</summary>

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
  <summary> 预先偷看缓存中还未读取的数据 </summary>

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
  <summary> 探索Explore类函数 </summary>

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


## 参考

- https://github.com/smallnest/ringbuffer
- https://github.com/Allenxuxu/ringbuffer

## 附录

欢迎PR
