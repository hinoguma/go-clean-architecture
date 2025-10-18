package service

import (
	"app/threelayeredarchitecture/appinfraadapterlayer"
	"app/threelayeredarchitecture/applogiclayer/usecase/domain/model"
)

type BankAccountLockServiceIF interface {
	Lock(bankAccountID string) (model.BankAccount, error)
	Unlock(bankAccountID string) (model.BankAccount, error)
}

type BankAccountLockService struct {
	bankAccountRepository appinfraadapterlayer.BankAccountRepositoryAdapterIF
}

func NewBankAccountLockService(
	bankAccountRepository appinfraadapterlayer.BankAccountRepositoryAdapterIF,
) *BankAccountLockService {
	return &BankAccountLockService{
		bankAccountRepository: bankAccountRepository,
	}
}

func (s *BankAccountLockService) Lock(bankAccountID string) (model.BankAccount, error) {
	panic("not implemented")
}

func (s *BankAccountLockService) Unlock(bankAccountID string) (model.BankAccount, error) {
	panic("not implemented")
}
