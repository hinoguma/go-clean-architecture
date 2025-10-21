package service

import (
	"app/threelayeredarchitecture/appinfraadapterlayer"
	"app/threelayeredarchitecture/applogiclayer/usecase/domain/model"
	"context"
)

type BankAccountLockServiceIF interface {
	Lock(ctx context.Context, bankAccountID string) (model.BankAccount, error)
	Unlock(ctx context.Context, bankAccountID string) (model.BankAccount, error)
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

func (s *BankAccountLockService) Lock(ctx context.Context, bankAccountID string) (model.BankAccount, error) {
	panic("not implemented")
}

func (s *BankAccountLockService) Unlock(ctx context.Context, bankAccountID string) (model.BankAccount, error) {
	panic("not implemented")
}
