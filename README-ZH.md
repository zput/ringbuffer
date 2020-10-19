# 多功能的环形缓存

[![Github Actions](https://github.com/Allenxuxu/gev/workflows/CI/badge.svg)](https://github.com/Allenxuxu/gev/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/Allenxuxu/gev)](https://goreportcard.com/report/github.com/Allenxuxu/gev)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/a2a55fe9c0c443e198f588a6c8026cd0)](https://www.codacy.com/manual/Allenxuxu/gev?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=Allenxuxu/gev&amp;utm_campaign=Badge_Grade)
[![GoDoc](https://godoc.org/github.com/Allenxuxu/gev?status.svg)](https://godoc.org/github.com/Allenxuxu/gev)
[![LICENSE](https://img.shields.io/badge/LICENSE-MIT-blue)](https://github.com/Allenxuxu/gev/blob/master/LICENSE)
[![Code Size](https://img.shields.io/github/languages/code-size/Allenxuxu/gev.svg?style=flat)](https://img.shields.io/github/languages/code-size/Allenxuxu/gev.svg?style=flat)

#### 中文 | [English](README.md)

环形缓存：
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

> 测试环境 Mac 

### 吞吐量测试

无锁

![image](benchmarks/out/unlock.jpg)

有锁

![image](benchmarks/out/lock.jpg)

</details>

## 安装

```bash
go get -u github.com/zput/ringbuffer
```

## 示例

<details>
  <summary> 创建有锁/无锁ringbuffer对象</summary>

```go

```

</details>

<details>
  <summary> 预先偷看缓存中还未读取的数据 </summary>

```go

```

</details>

<details>
  <summary> 探索Explore类函数 </summary>

```go
```

</details>


## 参考

- https://github.com/smallnest/ringbuffer
- https://github.com/Allenxuxu/ringbuffer

## 附录

欢迎PR
