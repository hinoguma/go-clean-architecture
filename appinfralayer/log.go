package appinfralayer

import (
	"app/crosscutting/infrainterface"
	"context"
	"fmt"
)

type cloudWatchLogger struct {
}

func (c cloudWatchLogger) Info(ctx context.Context, req infrainterface.LogRequest) {
	logjson := req.LogJson("INFO")
	fmt.Println(logjson.JsonString(), "from cloudWatchLogger")
}

func (c cloudWatchLogger) Error(ctx context.Context, req infrainterface.LogRequest) {
	logjson := req.LogJson("ERROR")
	fmt.Println(logjson.JsonString(), "from cloudWatchLogger")
}

func (c cloudWatchLogger) Debug(ctx context.Context, req infrainterface.LogRequest) {
	logjson := req.LogJson("DEBUG")
	fmt.Println(logjson.JsonString(), "from cloudWatchLogger")
}

func NewCloudWatchLogger() infrainterface.Logger {
	return cloudWatchLogger{}
}
