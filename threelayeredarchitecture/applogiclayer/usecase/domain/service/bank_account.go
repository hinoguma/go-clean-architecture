package service

import (
	"app/threelayeredarchitecture/appinfraadapterlayer"
	"app/threelayeredarchitecture/applogiclayer/usecase/domain/model"
	"context"
)

type BankAccountLockServiceIF interface {
	Lock(ctx context.Context, bankAccountID string, tx appinfraadapterlayer.Transaction) (model.BankAccount, error)
	//Unlock(ctx context.Context, bankAccountID string, tx appinfraadapterlayer.Transaction) (model.BankAccount, error)
}

type BankAccountLockService struct {
	bankAccountRepository appinfraadapterlayer.BankAccountRepositoryAdapterIF
}

func NewBankAccountLockService(
	bankAccountRepository appinfraadapterlayer.BankAccountRepositoryAdapterIF,
) BankAccountLockServiceIF {
	return &BankAccountLockService{
		bankAccountRepository: bankAccountRepository,
	}
}

func (s *BankAccountLockService) Lock(ctx context.Context, bankAccountID string, tx appinfraadapterlayer.Transaction) (model.BankAccount, error) {

	itemDto, err := s.bankAccountRepository.Lock(ctx, bankAccountID, tx)
	if err != nil {
		return model.BankAccount{}, err
	}
	item := model.BankAccount{}
	item.SetFromDTO(itemDto)
	return item, nil
}

//
//func (s *BankAccountLockService) Unlock(ctx context.Context, bankAccountID string, exTx appinfraadapterlayer.Transaction) (model.BankAccount, error) {
//
//}
