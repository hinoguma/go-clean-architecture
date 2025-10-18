package service

import (
	"app/threelayeredarchitecture/appinfraadapterlayer"
	"app/threelayeredarchitecture/applogiclayer/usecase/domain/model"
	"errors"
	"fmt"
	"time"
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

func (service *TransferService) Execute(req model.TransferRequest) (model.TransferResult, error) {

	// rollback mechanism is omitted for simplicity
	// my current strategy would be... Maiking WAL log struct and if some error happens, we can rollback to the previous state
	result := model.TransferResult{}

	// get bank account and user
	fromAccountDTO, err := service.bankAccountRepository.Get(req.FromBankAccountID)
	if err != nil {
		return result, err
	}
	toAccountDTO, err := service.bankAccountRepository.Get(req.ToBankAccountID)
	if err != nil {
		return result, err
	}
	fromAccount := model.ConvertDTOToBankAccount(fromAccountDTO)
	toAccount := model.ConvertDTOToBankAccount(toAccountDTO)

	// lock accounts
	_, err = service.bankAccountLockService.Lock(fromAccount.ID)
	if err != nil {
		return result, err
	}
	defer service.bankAccountLockService.Unlock(fromAccount.ID)
	_, err = service.bankAccountLockService.Lock(toAccount.ID)
	if err != nil {
		return result, err
	}
	defer service.bankAccountLockService.Unlock(toAccount.ID)

	// create transaction record
	record := model.NewTransferRecord(
		"uuid v4", fromAccount.ID, toAccount.ID, req.Money, time.Now(),
	)
	record.SetTransferStatus(model.TransferStatusInCheckBalance)
	err = service.transactionRecordRepository.Create(record.DTO())
	if err != nil {
		return result, err
	}

	// check balance
	if fromAccount.Balance.IsLessThan(req.Money) {
		result.TransactionRecord.SetTransferStatus(model.TransferStatusInCheckBalanceFailed)
		err2 := service.transactionRecordRepository.Update(record.DTO())
		if err2 != nil {
			return result, fmt.Errorf("insufficient balance: %w", err2)
		}
		return result, errors.New("insufficient balance")
	}

	// withdraw from the fromAccount
	result.TransactionRecord.SetTransferStatus(
		model.TransferStatusInTransfer,
	)
	result.TransactionRecord.UpdatedAt = time.Now()
	service.transactionRecordRepository.Update(result.TransactionRecord.DTO())
	fromAccount.Withdraw(req.Money, time.Now())
	err = service.bankAccountRepository.Update(fromAccount.DTO())
	if err != nil {
		wrappedErr := service.rollbackTransfer(
			req, fromAccount, toAccount, result.TransactionRecord, err,
		)
		return result, wrappedErr
	}

	// deposit to the toAccount
	toAccount.Deposit(req.Money, time.Now())
	err = service.bankAccountRepository.Update(toAccount.DTO())
	if err != nil {
		wrappedErr := service.rollbackTransfer(
			req, fromAccount, toAccount, result.TransactionRecord, err,
		)
		return result, wrappedErr
	}

	// set transfer completed
	result.TransactionRecord.SetTransferStatus(model.TransferStatusCompleted)
	result.TransactionRecord.UpdatedAt = time.Now()
	err = service.transactionRecordRepository.Update(result.TransactionRecord.DTO())
	if err != nil {
		wrappedErr := service.rollbackTransfer(
			req, fromAccount, toAccount, result.TransactionRecord, err,
		)
		return result, wrappedErr
	}

	// success
	return result, nil
}

func (service *TransferService) rollbackTransfer(
	req model.TransferRequest,
	fromAccount model.BankAccount,
	toAccount model.BankAccount,
	record model.TransactionRecord,
	err error,
) error {
	// omitted for simplicity
	return nil
}
