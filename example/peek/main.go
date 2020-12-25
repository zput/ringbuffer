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

	var isUsingExplore = false

	print(buffer.Peek(7, isUsingExplore))

	print(buffer.PeekAll(isUsingExplore))

	fmt.Println(buffer.PeekUint8(isUsingExplore))

	fmt.Println(buffer.PeekUint16(isUsingExplore))

	fmt.Println(buffer.PeekUint32(isUsingExplore))

	fmt.Println(buffer.PeekUint64(isUsingExplore))
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
