package redislock

import (
	"context"
	"time"

	"github.com/WeiJiadong/redislock/internal"
	"github.com/go-redis/redis/v9"
)

// Mutex 互斥锁定义
type Mutex struct {
	internal.LockInfo
}

// Lock 互斥锁加锁
func (l *Mutex) Lock(ctx context.Context, client *redis.Client) (bool, error) {
	return client.SetNX(ctx, l.Key, l.Val, l.Lease).Result()
}

// Unlock 互斥锁解锁
func (l *Mutex) Unlock(ctx context.Context, client *redis.Client) error {
	return internal.ExecLuaScript(ctx, client, l.Key, internal.UnLockScript, l.Val)
}

// Refresh 刷新互斥锁租约
func (l *Mutex) Refresh(ctx context.Context, client *redis.Client) error {
	return internal.ExecLuaScript(ctx, client, l.Key, internal.RefreshScript, l.Val, l.Lease/time.Second)
}

// LockOrRefresh 互斥锁加锁或者刷新租约
func (l *Mutex) LockOrRefresh(ctx context.Context, client *redis.Client) error {
	return internal.ExecLuaScript(ctx, client, l.Key, internal.LockOrRefreshScript, l.Val, l.Lease/time.Second)
}

// NewMutex 互斥锁构造函数
func NewMutex(Key, Val string, Lease time.Duration) *Mutex {
	return &Mutex{
		LockInfo: internal.LockInfo{
			Key:   Key,
			Val:   Val,
			Lease: Lease,
		}}
}
