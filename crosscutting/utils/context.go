package utils

import "context"

const CtxRequestIDKey = "requestId"

func GetRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v, ok := ctx.Value(CtxRequestIDKey).(string); ok {
		return v
	}
	return ""
}

func SetRequestID(ctx context.Context, requestID string) context.Context {
	if ctx == nil {
		return ctx
	}
	return context.WithValue(ctx, CtxRequestIDKey, requestID)
}
