#  multi-function ringbuffer

[![LICENSE](https://img.shields.io/badge/LICENSE-MIT-blue)](https://github.com/zput/ringbuffer/blob/master/LICENSE)
[![Github Actions](https://github.com/zput/ringbuffer/workflows/CI/badge.svg)](https://github.com/zput/ringbuffer/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/zput/ringbuffer)](https://goreportcard.com/report/github.com/zput/ringbuffer)
[![GoDoc](https://godoc.org/github.com/zput/ringbuffer?status.svg)](https://godoc.org/github.com/zput/ringbuffer)

#### [中文](README-ZH.md) | English

- Control whether locking(thread safe) or unlocking(single thread; fast) is required via parameters
- Automatic expansion of the circular buffer implementation
- Pre-read the data in the cache by exploring(ExploreBegin()----ExploreRead()/ExploreSize()----ExploreCommit()/ExploreBreak())


## example

```go
//TODO

```

## reference

- https://github.com/smallnest/ringbuffer
- https://github.com/Allenxuxu/ringbuffer

## appendix

welcome pr