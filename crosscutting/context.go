package crosscutting

import "context"

const CtxRequestIDKey = "requestId"

func GetRequestID(ctx context.Context) string {
	if v, ok := ctx.Value(CtxRequestIDKey).(string); ok {
		return v
	}
	return ""
}

func SetRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, CtxRequestIDKey, requestID)
}
