// Package rwlock 读写锁实现
package internal

import (
	"context"
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
)

func TestNewRWMutex(t *testing.T) {
	ctx := context.TODO()
	cli := getClient(ctx)
	type args struct {
		key   string
		val   string
		lease time.Duration
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
				key:   "key1",
				val:   "val1",
				lease: 2 * time.Second,
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
			// 先清理，防止其他类别key存在导致出错
			cli.Del(ctx, tt.args.key)
			locker := NewRWMutex(tt.args.key, tt.args.val, tt.args.lease)
			// 先加写锁
			err := locker.Lock(ctx, cli)
			assert.Equal(t, err, tt.wants.err)
			// check 是不是租约的时间
			ttl, err := cli.TTL(ctx, tt.args.key).Result()
			assert.Equal(t, ttl, tt.wants.ttl)
			assert.Equal(t, err, tt.wants.err)
			// check 自动解锁是否成功
			time.Sleep(ttl + ttl/2)
			cnt, err := cli.Exists(ctx, tt.args.key).Result()
			assert.Equal(t, cnt, tt.wants.cnt)
			assert.Equal(t, err, tt.wants.err)
			// 加读锁
			assert.Equal(t, locker.RWLock(ctx, cli), tt.wants.err)
			// 加写锁
			assert.Equal(t, locker.Lock(ctx, cli), tt.wants.err)
			// check加锁是否成功，以及是不是租约时间
			ttl, err = cli.TTL(ctx, tt.args.key).Result()
			assert.Equal(t, ttl, tt.wants.ttl)
			assert.Equal(t, err, tt.wants.err)
			// 释放锁
			assert.Equal(t, locker.Unlock(ctx, cli), tt.wants.err)
			// check 解锁是否成功
			cnt, err = cli.Exists(ctx, tt.args.key).Result()
			assert.Equal(t, cnt, tt.wants.cnt)
			assert.Equal(t, err, tt.wants.err)
		})
	}
}
