package log

import (
	"context"
	"github.com/rafian-git/go-backend/pkg/appctx"
	"go.uber.org/zap"
	"testing"
)

func TestLogger_Info(t *testing.T) {
	type fields struct {
		Logger *zap.Logger
	}
	type args struct {
		ctx    context.Context
		msg    string
		fields []zap.Field
	}

	cx := context.Background()
	cx = WithTraceID(cx, "111111111")
	cx = appctx.WithUserId(cx, "123456")
	//cx = context.WithT(cx, RequestIDContextKey, )
	//cx = context.WithValue(cx, TraceIDContextKey, "2222222222")

	var param []zap.Field
	//param = append(param, zap.String("test", "123"))
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: args{msg: "test message", fields: param, ctx: cx},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			log := New().Named("test")
			log.Info(tt.args.ctx, tt.args.msg, zap.String("test", "9999"))
		})
	}
}
