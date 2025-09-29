package infrainterface

import (
	"context"
)

// implements in infrastructure layer
type Logger interface {
	Info(ctx context.Context, req LogRequest)
	Error(ctx context.Context, req LogRequest)
	Debug(ctx context.Context, req LogRequest)
}

type LogRequest interface {
	LogJson(level string) LogJson
}

type LogJson interface {
	JsonString() string
}
