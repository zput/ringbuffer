package ringbuffer

//
//import "fmt"
//
//func ExampleRingBuffer() {
//	rb := New(1024)
//	_, _ = rb.Write([]byte("abcd"))
//	fmt.Println(rb.Size())
//	fmt.Println(rb.free())
//	buf := make([]byte, 4)
//
//	_, _ = rb.Read(buf)
//	fmt.Println(string(buf))
//
//	rb.Write([]byte("1234567890"))
//	rb.ExploreRead(buf)
//	fmt.Println(string(buf))
//	fmt.Println(rb.Size())
//	fmt.Println(rb.ExploreSize())
//	rb.ExploreFlush()
//	fmt.Println(rb.Size())
//	fmt.Println(rb.ExploreSize())
//
//	rb.ExploreRead(buf)
//	fmt.Println(string(buf))
//	fmt.Println(rb.Size())
//	fmt.Println(rb.ExploreSize())
//	rb.ExploreRevert()
//	fmt.Println(rb.Size())
//	fmt.Println(rb.ExploreSize())
//	// Output: 4
//	// 1020
//	// abcd
//	// 1234
//	// 10
//	// 6
//	// 6
//	// 6
//	// 5678
//	// 6
//	// 2
//	// 6
//	// 6
//}
