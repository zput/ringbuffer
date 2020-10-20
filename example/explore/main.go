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
