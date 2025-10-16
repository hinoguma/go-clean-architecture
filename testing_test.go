package go_clean_architecture

import (
	"app/experimentarchitecture/crosscutting/logger"
	"app/experimentarchitecture/crosscutting/utils"
	"testing"
)

func TestA(t *testing.T) {

	logger.LogInfo(nil, utils.NewLogRequest("aho!").ToValue())

}
