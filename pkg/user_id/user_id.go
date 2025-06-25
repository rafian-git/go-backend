package user_id

import (
	"context"

	"github.com/rafian-git/go-backend/pkg/appctx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UserIdUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return handler(ctx, req)
	}

	if len(md[appctx.USER_ID_HEADER]) > 0 {
		userId := md[appctx.USER_ID_HEADER][0]
		ctx = appctx.WithUserId(ctx, userId)
	}

	if len(md[appctx.USER_NAME]) > 0 {
		userName := md[appctx.USER_NAME][0]
		ctx = appctx.SetUserName(ctx, userName)
	}

	if len(md[appctx.EMAIL_HEADER]) > 0 {
		email := md[appctx.EMAIL_HEADER][0]
		ctx = appctx.WithEmail(ctx, email)
	}

	if len(md[appctx.USER_PHONE]) > 0 {
		phone := md[appctx.USER_PHONE][0]
		ctx = appctx.SetUserPhone(ctx, phone)
	}

	if len(md[appctx.USER_CLIENT_CODE]) > 0 {
		clientCode := md[appctx.USER_CLIENT_CODE][0]
		ctx = appctx.SetClientCode(ctx, clientCode)
	}

	return handler(ctx, req)
}
