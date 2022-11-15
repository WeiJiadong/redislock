# redislock
[![Badge](https://img.shields.io/badge/link-996.icu-%23FF4D5B.svg?style=flat-square)](https://996.icu/#/zh_CN)
[![Go](https://github.com/WeiJiadong/redislock/workflows/Go/badge.svg?branch=master)](https://github.com/WeiJiadong/redislock/actions)
[![Go Report Card](https://img.shields.io/badge/go%20report-A+-brightgreen.svg?style=flat)](https://goreportcard.com/report/github.com/WeiJiadong/redislock)
[![Latest](https://img.shields.io/badge/latest-v1.1.2-blue.svg)](https://github.com/WeiJiadong/redislock/tree/v1.1.2)
[![codecov](https://codecov.io/gh/WeiJiadong/redislock/branch/master/graph/badge.svg?token=6RG0W91RF2)](https://codecov.io/gh/WeiJiadong/redislock)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
### 项目介绍
go实现的一个基于lua+redis的分布式锁，项目依赖`github.com/go-redis/redis`的client实现。

### 功能介绍
1. 支持互斥锁
2. 支持读写锁

### 使用示例
```go
    cli := redis.NewClient(&redis.Options {
		Addr: "127.0.0.1:6379",
	})
    locker := NewRWMutex("key", "val", time.Second())
    locker.Lock(ctx, cli)
    locker.Lock(ctx, cli)
```