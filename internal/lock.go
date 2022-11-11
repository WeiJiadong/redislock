package internal

import (
	"context"
	"time"

	"github.com/go-redis/redis/v9"
)

// Mutex 互斥锁定义
type Mutex struct {
	lockInfo
}

// Lock 互斥锁加锁
func (l *Mutex) Lock(ctx context.Context, client *redis.Client) (bool, error) {
	return client.SetNX(ctx, l.key, l.val, l.lease/time.Second).Result()
}

// Unlock 互斥锁解锁
func (l *Mutex) Unlock(ctx context.Context, client *redis.Client) error {
	return execLuaScript(ctx, client, l.key, unLockScript, l.val)
}

// Refresh 刷新互斥锁租约
func (l *Mutex) Refresh(ctx context.Context, client *redis.Client) error {
	return execLuaScript(ctx, client, l.key, refreshScript, l.val, l.lease/time.Second)
}

// LockOrRefresh 互斥锁加锁或者刷新租约
func (l *Mutex) LockOrRefresh(ctx context.Context, client *redis.Client) error {
	return execLuaScript(ctx, client, l.key, lockOrRefreshScript, l.val, l.lease/time.Second)
}

// NewMutex 互斥锁构造函数
func NewMutex(key, val string, lease time.Duration) *Mutex {
	return &Mutex{
		lockInfo: lockInfo{
			key:   key,
			val:   val,
			lease: lease,
		}}
}
