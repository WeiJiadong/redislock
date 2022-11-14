package redislock

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/WeiJiadong/redislock/internal"
	"github.com/go-redis/redis/v9"
	"gopkg.in/go-playground/assert.v1"
)

func getClient(ctx context.Context) *redis.Client {
	rds := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})

	_, err := rds.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	return rds
}

func TestNewMutex(t *testing.T) {
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
			// 先清理，防止其他类别Key存在导致出错
			cli.Del(ctx, tt.args.Key)
			locker := NewMutex(tt.args.Key, tt.args.Val, tt.args.Lease)
			// 先加锁
			assert.Equal(t, locker.Lock(ctx, cli), tt.wants.err)
			// check 是不是租约的时间
			ttl, err := cli.TTL(ctx, tt.args.Key).Result()
			assert.Equal(t, ttl, tt.wants.ttl)
			assert.Equal(t, err, tt.wants.err)
			// 模拟逻辑阻塞，触发续约
			time.Sleep(tt.args.Lease / 2)
			assert.Equal(t, locker.Refresh(ctx, cli), tt.wants.err)
			// check 续约是否成功
			ttl, err = cli.TTL(ctx, tt.args.Key).Result()
			assert.Equal(t, ttl, tt.wants.ttl)
			assert.Equal(t, err, tt.wants.err)
			// check 自动解锁是否成功
			time.Sleep(ttl + ttl/2)
			cnt, err := cli.Exists(ctx, tt.args.Key).Result()
			assert.Equal(t, cnt, tt.wants.cnt)
			assert.Equal(t, err, tt.wants.err)
			// 通过LockOrRefresh加锁
			assert.Equal(t, locker.LockOrRefresh(ctx, cli), tt.wants.err)
			// check加锁是否成功，以及是不是租约时间
			ttl, err = cli.TTL(ctx, tt.args.Key).Result()
			assert.Equal(t, ttl, tt.wants.ttl)
			assert.Equal(t, err, tt.wants.err)
			// 模拟逻辑阻塞，触发续约
			time.Sleep(tt.args.Lease / 2)
			assert.Equal(t, locker.LockOrRefresh(ctx, cli), tt.wants.err)
			// check 续约是否成功
			ttl, err = cli.TTL(ctx, tt.args.Key).Result()
			assert.Equal(t, ttl, tt.wants.ttl)
			assert.Equal(t, err, tt.wants.err)
			// 释放锁
			assert.Equal(t, locker.Unlock(ctx, cli), tt.wants.err)
			// check 解锁是否成功
			cnt, err = cli.Exists(ctx, tt.args.Key).Result()
			assert.Equal(t, cnt, tt.wants.cnt)
			assert.Equal(t, err, tt.wants.err)
			// 加两次锁
			assert.Equal(t, locker.Lock(ctx, cli), tt.wants.err)
			assert.Equal(t, locker.Lock(ctx, cli), errors.New(internal.ErrLockFailed))
			// 释放锁，对不存在的锁进行续约
			assert.Equal(t, locker.Unlock(ctx, cli), tt.wants.err)
			assert.Equal(t, locker.Refresh(ctx, cli), errors.New(internal.ErrLockNotExist))
		})
	}
}
