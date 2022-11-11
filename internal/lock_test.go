package internal

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/magiconair/properties/assert"
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
			locker := NewMutex(tt.args.key, tt.args.val, tt.args.lease)
			// 先加锁
			ok, err := locker.Lock(ctx, cli)
			assert.Equal(t, ok, true)
			assert.Equal(t, err, tt.wants.err)
			// check 是不是租约的时间
			ttl, err := cli.TTL(ctx, tt.args.key).Result()
			assert.Equal(t, ttl, tt.wants.ttl)
			assert.Equal(t, err, tt.wants.err)
			// 模拟逻辑阻塞，触发续约
			time.Sleep(tt.args.lease / 2)
			assert.Equal(t, locker.Refresh(ctx, cli), tt.wants.err)
			// check 续约是否成功
			ttl, err = cli.TTL(ctx, tt.args.key).Result()
			assert.Equal(t, ttl, tt.wants.ttl)
			assert.Equal(t, err, tt.wants.err)
			// check 自动解锁是否成功
			time.Sleep(ttl + ttl/2)
			cnt, err := cli.Exists(ctx, tt.args.key).Result()
			assert.Equal(t, cnt, tt.wants.cnt)
			assert.Equal(t, err, tt.wants.err)
			// 通过LockOrRefresh加锁
			assert.Equal(t, locker.LockOrRefresh(ctx, cli), tt.wants.err)
			// check加锁是否成功，以及是不是租约时间
			ttl, err = cli.TTL(ctx, tt.args.key).Result()
			assert.Equal(t, ttl, tt.wants.ttl)
			assert.Equal(t, err, tt.wants.err)
			// 模拟逻辑阻塞，触发续约
			time.Sleep(tt.args.lease / 2)
			assert.Equal(t, locker.LockOrRefresh(ctx, cli), tt.wants.err)
			// check 续约是否成功
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
