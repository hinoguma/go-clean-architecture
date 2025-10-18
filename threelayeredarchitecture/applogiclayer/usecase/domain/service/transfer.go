package service

import (
	"app/threelayeredarchitecture/appinfraadapterlayer"
	"app/threelayeredarchitecture/applogiclayer/usecase/domain/model"
)

type TransferServiceIF interface {
	Execute(req model.TransferRequest) (model.TransferResult, error)
}

type TransferService struct {
	bankAccountLockService      BankAccountLockServiceIF
	bankCustomerRepository      appinfraadapterlayer.BankCustomerRepositoryAdapterIF
	bankAccountRepository       appinfraadapterlayer.BankAccountRepositoryAdapterIF
	transactionRecordRepository appinfraadapterlayer.TransactionRecordRepositoryAdapterIF
}

func NewTransferService(
	bankAccountLockService BankAccountLockServiceIF,
	bankCustomerRepository appinfraadapterlayer.BankCustomerRepositoryAdapterIF,
	bankAccountRepository appinfraadapterlayer.BankAccountRepositoryAdapterIF,
	transactionRecordRepository appinfraadapterlayer.TransactionRecordRepositoryAdapterIF,
) *TransferService {
	return &TransferService{
		bankAccountLockService:      bankAccountLockService,
		bankCustomerRepository:      bankCustomerRepository,
		bankAccountRepository:       bankAccountRepository,
		transactionRecordRepository: transactionRecordRepository,
	}
}

func (s *TransferService) Execute(req model.TransferRequest) (model.TransferResult, error) {
	panic("not implemented")
}
