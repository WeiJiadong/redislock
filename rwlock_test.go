package redislock

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/WeiJiadong/redislock/internal"
	"github.com/magiconair/properties/assert"
)

func TestNewRWMutex(t *testing.T) {
	ctx := context.TODO()
	cli := getClient(ctx)
	type args struct {
		Key   string
		Val   string
		Lease time.Duration
	}
	type wants struct {
		ttl time.Duration
		cnt int64
		err error
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name: "普通case:",
			args: args{
				Key:   "Key1",
				Val:   "Val1",
				Lease: 2 * time.Second,
			},
			wants: wants{
				ttl: 2 * time.Second,
				cnt: int64(0),
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			locker := NewRWMutex(tt.args.Key, tt.args.Val, tt.args.Lease)
			// 先清理，防止其他类别Key存在导致出错
			cli.Del(ctx, tt.args.Key)
			// 先加写锁
			assert.Equal(t, locker.WLock(ctx, cli), tt.wants.err)
			// check 是不是租约的时间
			ttl, err := cli.TTL(ctx, tt.args.Key).Result()
			assert.Equal(t, ttl, tt.wants.ttl)
			assert.Equal(t, err, tt.wants.err)
			// check 自动解锁是否成功
			time.Sleep(ttl + ttl/2)
			cnt, err := cli.Exists(ctx, tt.args.Key).Result()
			assert.Equal(t, cnt, tt.wants.cnt)
			assert.Equal(t, err, tt.wants.err)
			// 加读锁
			// 先清理，防止其他类别Key存在导致出错
			assert.Equal(t, locker.RLock(ctx, cli), tt.wants.err)
			// 加写锁
			assert.Equal(t, locker.WLock(ctx, cli), tt.wants.err)
			// check加锁是否成功，以及是不是租约时间
			ttl, err = cli.TTL(ctx, tt.args.Key).Result()
			assert.Equal(t, ttl, tt.wants.ttl)
			assert.Equal(t, err, tt.wants.err)
			// 释放写锁
			assert.Equal(t, locker.WUnlock(ctx, cli), tt.wants.err)
			// 释放读锁
			assert.Equal(t, locker.RUnlock(ctx, cli), tt.wants.err)
			// check 解锁是否成功
			cnt, err = cli.Exists(ctx, tt.args.Key).Result()
			assert.Equal(t, cnt, tt.wants.cnt)
			assert.Equal(t, err, tt.wants.err)
			// 重复加写锁
			assert.Equal(t, locker.WLock(ctx, cli), tt.wants.err)
			assert.Equal(t, locker.WLock(ctx, cli), errors.New(internal.ErrWLockConflict))
			// 加读锁
			assert.Equal(t, locker.RLock(ctx, cli), errors.New(internal.ErrHadWLock))
		})
	}
}
