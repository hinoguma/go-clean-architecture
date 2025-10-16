package infrainterface

import (
	"app/experimentarchitecture/crosscutting/utils"
	"context"
)

// implements in infra layer
type Logger interface {
	Info(ctx context.Context, req utils.LogRequest)
	Error(ctx context.Context, req utils.LogRequest)
	Debug(ctx context.Context, req utils.LogRequest)
}
