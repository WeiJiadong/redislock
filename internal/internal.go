// Package internal 依赖的内部实现
package internal

import (
	"context"
	_ "embed"
	"time"

	"github.com/go-redis/redis/v9"
)

// lockInfo 锁结构信息
type lockInfo struct {
	key   string
	val   string
	lease time.Duration
}

// lua脚本变量定义，直接将lua脚本加载进二进制
var (
	//go:embed script/lock_or_refresh.lua
	lockOrRefreshScript string
	//go:embed script/refresh.lua
	refreshScript string
	//go:embed script/unlock.lua
	unLockScript string
	//go:embed script/rlock.lua
	rLockScript string
	//go:embed script/runlock.lua
	rUnlockScript string
	//go:embed script/wlock.lua
	wLockScript string
	//go:embed script/wunlock.lua
	wUnlockScript string
)

// execLuaScript 执行lua脚本
func execLuaScript(ctx context.Context, client *redis.Client, key, script string, vals ...any) error {
	lua := redis.NewScript(script)
	_, err := lua.Run(ctx, client, []string{key}, vals).Result()
	return err
}
