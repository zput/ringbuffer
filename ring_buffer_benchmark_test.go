package ringbuffer

import (
	"runtime"
	"strings"
	"testing"
)

func BenchmarkRingBuffer_Sync_Unlock(b *testing.B) {
	rb := New(1024)
	data := []byte(strings.Repeat("a", 512))
	buf := make([]byte, 512)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = rb.Write(data)
		_, _ = rb.Read(buf)
	}
}

func BenchmarkRingBuffer_Sync_Lock(b *testing.B) {
	rb := New(1024, true)
	data := []byte(strings.Repeat("a", 512))
	buf := make([]byte, 512)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = rb.Write(data)
		_, _ = rb.Read(buf)
	}
}

func BenchmarkRingBuffer_AsyncWrite_Need_Lock(b *testing.B) {
	rb := New(1024, true)
	data := []byte(strings.Repeat("a", 512))
	buf := make([]byte, 512)
	var (
		errWrite error
		errWriteCount int

		errRead error
		errReadCount int
	)

	go func() {
		for {
			_, errWrite = rb.Write(data)
			if errWrite != nil {
				errWriteCount++
				runtime.Gosched()
			}
		}
	}()


	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, errRead = rb.Read(buf)
		if errRead != nil {
			errReadCount++
			runtime.Gosched()
		}
	}
	b.StopTimer()
	b.Logf("errReadCount[%d], errWriteCount[%d], b.N[%d]", errReadCount, errWriteCount, b.N)
	b.Log(errRead)
}

func BenchmarkRingBuffer_AsyncRead_Lock(b *testing.B) {
	rb := New(1024, true)
	data := []byte(strings.Repeat("a", 512))
	buf := make([]byte, 512)

	go func() {
		for {
			rb.Read(buf)
		}
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rb.Write(data)
	}
}

func BenchmarkRingBuffer_AsyncRead_Lock_Print_Log(b *testing.B) {
	rb := New(1024)
	data := []byte(strings.Repeat("a", 512))
	buf := make([]byte, 512)

	go func() {
		for {
			b.Log(rb.Read(buf))
		}
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.Log(rb.Write(data))
	}
}
