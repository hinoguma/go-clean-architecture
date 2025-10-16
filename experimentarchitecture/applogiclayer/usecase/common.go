package usecase

import (
	"app/experimentarchitecture/applogiclayer/domain/model"
	"app/experimentarchitecture/crosscutting/utils"
	"context"
)

type BaseUseCaseRequest struct {
	Ctx             context.Context
	RequestUserInfo RequestUserInfo
}

func (req BaseUseCaseRequest) BankAccountUserAuthRequest() model.BankUserAuthRequest {
	return model.BankUserAuthRequest{
		Ctx:   req.Ctx,
		Token: req.RequestUserInfo.AuthToken,
	}
}

type RequestUserInfo struct {
	IP        utils.IPAddress
	AuthToken string
}
