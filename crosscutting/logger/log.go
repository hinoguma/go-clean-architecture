package logger

import (
	"app/crosscutting/logger/infrainterface"
	"app/crosscutting/utils"
	"context"
	"fmt"
)

func SetGlobalLogger(logger infrainterface.Logger) {
	globalLogger = logger
}

var globalLogger infrainterface.Logger = NewStdLogger()

func LogInfo(ctx context.Context, req utils.LogRequest) {
	if globalLogger == nil {
		return
	}
	globalLogger.Info(ctx, req)
}

type StdLogger struct {
}

func NewStdLogger() infrainterface.Logger {
	return StdLogger{}
}

func (s StdLogger) Info(ctx context.Context, req utils.LogRequest) {
	logjson := req.LogJson()
	logjson.WithLevel(utils.LogLevelInfo).WithRequestIDByCtx(ctx)
	fmt.Println(logjson.JsonString())
}
func (s StdLogger) Error(ctx context.Context, req utils.LogRequest) {
	logjson := req.LogJson()
	logjson.WithLevel(utils.LogLevelError).WithRequestIDByCtx(ctx)
	fmt.Println(logjson.JsonString())
}

func (s StdLogger) Debug(ctx context.Context, req utils.LogRequest) {
	logjson := req.LogJson()
	logjson.WithLevel(utils.LogLevelDebug).WithRequestIDByCtx(ctx)
	fmt.Println(logjson.JsonString())
}
