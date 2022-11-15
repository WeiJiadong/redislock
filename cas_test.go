package redislock

import (
	"context"
	"errors"
	"testing"

	"github.com/WeiJiadong/redislock/internal"
	"gopkg.in/go-playground/assert.v1"
)

func TestNewCas(t *testing.T) {
	ctx := context.TODO()
	cli := getClient(ctx)
	type args struct {
		key    string
		val    string
		newVal *internal.CasVal
		oldVal *internal.CasVal
	}
	type wants struct {
		cnt int64
		err error
		val *internal.CasVal
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name: "普通case:",
			args: args{
				key: "Key1",
				val: "Val1",
				newVal: &internal.CasVal{
					Val: "new",
				},
				oldVal: &internal.CasVal{
					Val: "old",
				},
			},
			wants: wants{
				cnt: int64(0),
				err: nil,
				val: &internal.CasVal{
					Val: "new",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 先清理，防止其他类别Key存在导致出错
			cli.Del(ctx, tt.args.key)
			locker := NewCas(tt.args.key, tt.args.val)
			// 先执行cas
			assert.Equal(t, locker.Set(ctx, cli, tt.args.key, tt.args.oldVal, tt.args.newVal), tt.wants.err)
			assert.Equal(t, locker.Set(ctx, cli, tt.args.key, tt.args.oldVal, tt.args.newVal),
				errors.New(internal.ErrCasConflict))
			val, err := locker.Get(ctx, cli, tt.args.key)
			assert.Equal(t, val, tt.args.newVal)
			assert.Equal(t, err, tt.wants.err)
		})
	}
}
