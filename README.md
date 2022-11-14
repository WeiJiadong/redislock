# redislock
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