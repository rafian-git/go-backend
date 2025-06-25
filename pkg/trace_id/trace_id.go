package trace_id

import (
	"context"
	"gitlab.techetronventures.com/core/backend/pkg/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return handler(ctx, req)
	}

	if len(md[log.TraceIDContextKey]) > 0 {
		tracId := md[log.TraceIDContextKey][0]
		ctx = log.WithTraceID(ctx, tracId)
		ctx = metadata.AppendToOutgoingContext(ctx, log.TraceIDContextKey, log.TraceID(ctx))
	}
	return handler(ctx, req)
}
