package redislock

import (
	"context"
	"encoding/json"

	"github.com/WeiJiadong/redislock/internal"
	"github.com/go-redis/redis/v9"
)

type Cas struct {
	internal.LockInfo
}

// Set 对简单key-val的结构执行cas写入
func (c *Cas) Set(ctx context.Context, client *redis.Client, key string, oldVal, newVal *internal.CasVal) error {
	newVal.Version = oldVal.Version + 1
	oldCas, err := json.Marshal(oldVal)
	if err != nil {
		return err
	}
	newCas, err := json.Marshal(newVal)
	if err != nil {
		return err
	}
	return internal.ExecLuaScript(ctx, client, c.Key, internal.CasScript, oldCas, newCas)
}

// Get 对简单key-val的string结构执行cas读取
func (c *Cas) Get(ctx context.Context, client *redis.Client, key string) (*internal.CasVal, error) {
	val, err := client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	cVal := internal.CasVal{}
	err = json.Unmarshal(val, &cVal)
	if err != nil {
		return nil, err
	}
	return &cVal, err
}

// NewCas 互斥锁构造函数
func NewCas(Key, Val string) *Cas {
	return &Cas{
		LockInfo: internal.LockInfo{
			Key: Key,
		}}
}
