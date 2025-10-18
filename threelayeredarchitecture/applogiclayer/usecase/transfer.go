package usecase

import (
	"app/threelayeredarchitecture/appinfraadapterlayer"
	"app/threelayeredarchitecture/applogiclayer/usecase/domain/model"
	"app/threelayeredarchitecture/applogiclayer/usecase/domain/service"
)

type TransferUsecaseRequest struct {
	FromBankAccountID string
	ToBankAccountID   string
	Money             model.Money
}

type TransferUsecaseResponse struct {
	TransactionID string
}

type TransferUsecaseIF interface {
	Execute(req TransferUsecaseRequest) (TransferUsecaseResponse, error)
}

type TransferUsecase struct {
	transferService        service.TransferServiceIF
	bankAccountRepository  appinfraadapterlayer.BankAccountRepositoryAdapterIF
	bankCustomerRepository appinfraadapterlayer.BankCustomerRepositoryAdapterIF
}

func (tu *TransferUsecase) Execute(req TransferUsecaseRequest) (TransferUsecaseResponse, error) {

	// validate request

	// get bank account and user

	// lock accounts

	// check balance

	// perform transfer

	// create transaction record

	// persist changes

}
