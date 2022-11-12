package redislock

import (
	"context"
	"time"

	"github.com/WeiJiadong/redislock/internal"
	"github.com/go-redis/redis/v9"
)

// RWMutex 读写锁定义
type RWMutex struct {
	internal.LockInfo
}

// RWLock 读写锁加写锁
func (rw *RWMutex) WLock(ctx context.Context, client *redis.Client) error {
	return internal.ExecLuaScript(ctx, client, rw.Key, internal.WLockScript, rw.Val, rw.Lease/time.Second)
}

// RWLock 读写锁加读锁
func (rw *RWMutex) RLock(ctx context.Context, client *redis.Client) error {
	return internal.ExecLuaScript(ctx, client, rw.Key, internal.RLockScript, rw.Val, rw.Lease/time.Second)
}

// Unlock 读写锁解写锁
func (rw *RWMutex) WUnlock(ctx context.Context, client *redis.Client) error {
	return internal.ExecLuaScript(ctx, client, rw.Key, internal.WUnlockScript, rw.Val)
}

// Unlock 读写锁解读锁
func (rw *RWMutex) RUnlock(ctx context.Context, client *redis.Client) error {
	return internal.ExecLuaScript(ctx, client, rw.Key, internal.RUnlockScript, rw.Val)
}

// NewRWMutex 读写锁构造函数
func NewRWMutex(Key, Val string, Lease time.Duration) *RWMutex {
	return &RWMutex{
		LockInfo: internal.LockInfo{
			Key:   Key,
			Val:   Val,
			Lease: Lease,
		}}
}
