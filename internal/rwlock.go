// Package rwlock 读写锁实现
package internal

import (
	"context"
	"time"

	"github.com/go-redis/redis/v9"
)

// RWMutex 读写锁定义
type RWMutex struct {
	lockInfo
}

// RWLock 读写锁加写锁
func (rw *RWMutex) WLock(ctx context.Context, client *redis.Client) error {
	return execLuaScript(ctx, client, rw.key, wLockScript, rw.val, rw.lease/time.Second)
}

// RWLock 读写锁加读锁
func (rw *RWMutex) RLock(ctx context.Context, client *redis.Client) error {
	return execLuaScript(ctx, client, rw.key, rLockScript, rw.val, rw.lease/time.Second)
}

// Unlock 读写锁解写锁
func (rw *RWMutex) WUnlock(ctx context.Context, client *redis.Client) error {
	return execLuaScript(ctx, client, rw.key, wUnlockScript, rw.val)
}

// Unlock 读写锁解读锁
func (rw *RWMutex) RUnlock(ctx context.Context, client *redis.Client) error {
	return execLuaScript(ctx, client, rw.key, rUnlockScript, rw.val)
}

// NewRWMutex 读写锁构造函数
func NewRWMutex(key, val string, lease time.Duration) *RWMutex {
	return &RWMutex{
		lockInfo: lockInfo{
			key:   key,
			val:   val,
			lease: lease,
		}}
}
