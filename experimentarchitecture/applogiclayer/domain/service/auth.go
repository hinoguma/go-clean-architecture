package service

import (
	"app/experimentarchitecture/applogiclayer/domain/model"
	"app/experimentarchitecture/crosscutting/utils"
)

type BankUserAuthService interface {
	Authenticate(req model.BankUserAuthRequest) (model.BankUser, error)
}

type bankAccountUserAuthService struct {
}

func NewBankUserAuthService() BankUserAuthService {
	return &bankAccountUserAuthService{}
}

func (s *bankAccountUserAuthService) Authenticate(req model.BankUserAuthRequest) (model.BankUser, error) {
	var err error
	if req.Token == "" {
		err = utils.NewAuthorizeFailedError("missing auth token")
		return model.BankUser{}, err
	}

	return model.BankUser{}, nil
}
