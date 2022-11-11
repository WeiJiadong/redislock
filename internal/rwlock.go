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
func (l *RWMutex) Lock(ctx context.Context, client *redis.Client) error {
	return execLuaScript(ctx, client, l.key, rWLockScript, "w", l.lease/time.Second, l.val)
}

// RWLock 读写锁加读锁
func (l *RWMutex) RWLock(ctx context.Context, client *redis.Client) error {
	return execLuaScript(ctx, client, l.key, rWLockScript, "r", l.lease/time.Second, l.val)
}

// Unlock 读写锁解锁
func (l *RWMutex) Unlock(ctx context.Context, client *redis.Client) error {
	return execLuaScript(ctx, client, l.key, rWUnlockScript, l.val)
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
