package appctx

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

// context keys
type (
	clientIPContextKey  struct{} // client IP address
	userAgentContextKey struct{} // client user agent from header
)

// USER_ID_HEADER
var USER_ID_HEADER = "stockx-user-id"
var EMAIL_HEADER = "stockx-email"
var ORIGIN = "origin"
var USER_NAME = "stockx-user-name"
var USER_PHONE = "stockx-user-phone"
var USER_CLIENT_CODE = "stockx-user-client-code"
var USER_META_ID = "stockx-user-meta-id"

//
// ClientIP
//

// ClientIP from context.
func ClientIP(ctx context.Context) (clientIP string) {
	clientIP, _ = ctx.Value(clientIPContextKey{}).(string)
	return
}

// WithClientIP sets given client IP to given context.
func WithClientIP(ctx context.Context, clientIP string) context.Context {
	return context.WithValue(ctx, clientIPContextKey{}, clientIP)
}

// UserAgent from context.
func UserAgent(ctx context.Context) (userAgent string) {
	userAgent, _ = ctx.Value(userAgentContextKey{}).(string)
	return
}

// WithUserAgent sets given useragent to given context.
func WithUserAgent(ctx context.Context, clientIP string) context.Context {
	return context.WithValue(ctx, userAgentContextKey{}, clientIP)
}

func GetUserID(ctx context.Context) (int64, error) {
	userId, ok := ctx.Value(strings.ToLower(USER_ID_HEADER)).(string)
	if !ok {
		return 0, fmt.Errorf("invalid user_id in context")
	}

	u, err := strconv.Atoi(userId)
	if err != nil {
		return 0, err
	}

	return int64(u), nil
}

func GetEmail(ctx context.Context) (string, error) {
	email, ok := ctx.Value(EMAIL_HEADER).(string)
	if !ok {
		return "", fmt.Errorf("email not present in context")
	}
	return email, nil
}

func WithEmail(ctx context.Context, email string) context.Context {
	return context.WithValue(ctx, EMAIL_HEADER, email)
}

// WithClientIP sets given useragent to given context.
func WithUserId(ctx context.Context, userId string) context.Context {
	return context.WithValue(ctx, USER_ID_HEADER, userId)
}

func GetMetaId(ctx context.Context) (int64, error) {
	meta_id, ok := ctx.Value(strings.ToLower(USER_META_ID)).(string)
	if !ok {
		return 0, fmt.Errorf("invalid meta_id in context")
	}

	u, err := strconv.Atoi(meta_id)
	if err != nil {
		return 0, err
	}

	return int64(u), nil
}

func SetMetaId(ctx context.Context, metaId string) context.Context {
	return context.WithValue(ctx, USER_META_ID, metaId)
}

// SetUserName sets given username to given context.
func SetUserName(ctx context.Context, userName string) context.Context {
	return context.WithValue(ctx, USER_NAME, userName)
}

func GetUserName(ctx context.Context) (string, error) {
	userName, ok := ctx.Value(strings.ToLower(USER_NAME)).(string)
	if !ok {
		return "", fmt.Errorf("phone is not present in context")
	}

	return userName, nil
}

func SetUserPhone(ctx context.Context, phone string) context.Context {
	return context.WithValue(ctx, USER_PHONE, phone)
}

func GetUserPhone(ctx context.Context) (string, error) {
	userName, ok := ctx.Value(strings.ToLower(USER_PHONE)).(string)
	if !ok {
		return "", fmt.Errorf("phone is not present in context")
	}

	return userName, nil
}

func SetClientCode(ctx context.Context, clientCode string) context.Context {
	return context.WithValue(ctx, USER_CLIENT_CODE, clientCode)
}

func GetClientCode(ctx context.Context) (string, error) {
	clientCode, ok := ctx.Value(strings.ToLower(USER_CLIENT_CODE)).(string)
	if !ok {
		return "", fmt.Errorf("client code is not present in context")
	}

	return clientCode, nil
}
