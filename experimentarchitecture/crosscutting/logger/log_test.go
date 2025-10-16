package logger

import (
	"app/experimentarchitecture/crosscutting/utils"
	"context"
	"testing"
)

func TestLog(t *testing.T) {

	ctx := context.Background()

	LogInfo(ctx, utils.NewLogRequest("aho!").ToValue())
}
