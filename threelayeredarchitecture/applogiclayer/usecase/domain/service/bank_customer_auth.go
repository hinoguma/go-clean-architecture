package service

import (
	"app/threelayeredarchitecture/appinfraadapterlayer"
	"app/threelayeredarchitecture/applogiclayer/usecase/domain/model"
	"context"
)

type BankCustomerAuthServiceIF interface {
	Authenticate(ctx context.Context, req model.BankCustomerAuthRequest) model.BankCustomerAuthResult
}

func NewBankCustomerAuthService(
	bankCustomerRepository appinfraadapterlayer.BankCustomerRepositoryAdapterIF,
) BankCustomerAuthServiceIF {

	return bankCustomerAuthService{
		bankCustomerRepository: bankCustomerRepository,
	}
}

type bankCustomerAuthService struct {
	bankCustomerRepository appinfraadapterlayer.BankCustomerRepositoryAdapterIF
}

func (service bankCustomerAuthService) Authenticate(ctx context.Context, req model.BankCustomerAuthRequest) model.BankCustomerAuthResult {
	adapterRes := service.bankCustomerRepository.Authenticate(
		ctx, appinfraadapterlayer.AuthenticateRequestDTO{
			Token: req.Token, Timestamp: req.Timestamp,
		},
	)
	if adapterRes.Success {
		customer := model.BankCustomer{}
		customer.SetFromDTO(adapterRes.Customer)
		return model.BankCustomerAuthResult{
			Customer:    customer,
			ErrorReason: model.AuthErrorReason{},
			Err:         nil,
		}
	}

	// error
	return model.BankCustomerAuthResult{
		Customer: model.BankCustomer{},
		ErrorReason: model.AuthErrorReason{
			TokenExpired:  adapterRes.ErrorReason.TokenExpired,
			InvalidToken:  adapterRes.ErrorReason.InvalidToken,
			InternalError: adapterRes.ErrorReason.InternalError,
		},
		Err: adapterRes.Err,
	}
}
