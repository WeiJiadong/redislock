// Package internal 依赖的内部实现
package internal

import (
	"context"
	_ "embed"
	"errors"
	"time"

	"github.com/go-redis/redis/v9"
)

// LockInfo 锁结构信息
type LockInfo struct {
	Key   string
	Val   string
	Lease time.Duration
}

type CasVal struct {
	Val     any
	Version int64
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
	//go:embed script/cas.lua
	CasScript string
)

// ExecLuaScript 执行lua脚本
func ExecLuaScript(ctx context.Context, client *redis.Client, Key, script string, Vals ...any) error {
	lua := redis.NewScript(script)
	val, err := lua.Eval(ctx, client, []string{Key}, Vals...).Result()
	if err != nil {
		return err
	}
	msg, ok := val.(string)
	if !ok {
		return errors.New(ErrTypeNotMatch)
	}
	if msg != OK {
		return errors.New(msg)
	}
	return nil
}
