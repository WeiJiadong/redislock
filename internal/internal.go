// Package internal 依赖的内部实现
package internal

import (
	"context"
	_ "embed"
	"time"

	"github.com/go-redis/redis/v9"
)

// LockInfo 锁结构信息
type LockInfo struct {
	Key   string
	Val   string
	Lease time.Duration
}

// lua脚本变量定义，直接将lua脚本加载进二进制
var (
	//go:embed script/lock_or_refresh.lua
	LockOrRefreshScript string
	//go:embed script/refresh.lua
	RefreshScript string
	//go:embed script/unlock.lua
	UnLockScript string
	//go:embed script/rlock.lua
	RLockScript string
	//go:embed script/runlock.lua
	RUnlockScript string
	//go:embed script/wlock.lua
	WLockScript string
	//go:embed script/wunlock.lua
	WUnlockScript string
)

// ExecLuaScript 执行lua脚本
func ExecLuaScript(ctx context.Context, client *redis.Client, Key, script string, Vals ...any) error {
	lua := redis.NewScript(script)
	_, err := lua.Run(ctx, client, []string{Key}, Vals).Result()
	return err
}
