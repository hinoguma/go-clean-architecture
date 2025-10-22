package service

import (
	"app/threelayeredarchitecture/appinfraadapterlayer"
	"app/threelayeredarchitecture/applogiclayer/usecase/domain/model"
	"context"
	"errors"
	"time"
)

type TransferServiceIF interface {
	Execute(ctx context.Context, req model.TransferRequest) model.TransferResult
}

type TransferService struct {
	bankAccountLockService      BankAccountLockServiceIF
	transactionManager          appinfraadapterlayer.TransactionManagerAdapterIF
	bankCustomerRepository      appinfraadapterlayer.BankCustomerRepositoryAdapterIF
	bankAccountRepository       appinfraadapterlayer.BankAccountRepositoryAdapterIF
	transactionRecordRepository appinfraadapterlayer.TransactionRecordRepositoryAdapterIF
}

func NewTransferService(
	bankAccountLockService BankAccountLockServiceIF,
	transactionManager appinfraadapterlayer.TransactionManagerAdapterIF,
	bankCustomerRepository appinfraadapterlayer.BankCustomerRepositoryAdapterIF,
	bankAccountRepository appinfraadapterlayer.BankAccountRepositoryAdapterIF,
	transactionRecordRepository appinfraadapterlayer.TransactionRecordRepositoryAdapterIF,
) *TransferService {
	return &TransferService{
		bankAccountLockService:      bankAccountLockService,
		transactionManager:          transactionManager,
		bankCustomerRepository:      bankCustomerRepository,
		bankAccountRepository:       bankAccountRepository,
		transactionRecordRepository: transactionRecordRepository,
	}
}

func (service *TransferService) Execute(ctx context.Context, req model.TransferRequest) model.TransferResult {

	// rollback mechanism is omitted for simplicity
	// my current strategy would be... Maiking WAL log struct and if some error happens, we can rollback to the previous state
	result := model.TransferResult{}
	// Transaction Begins
	txBeginReq := appinfraadapterlayer.NewBeginTransactionRequest()
	tx, err := service.transactionManager.Begin(ctx, txBeginReq)
	if err != nil {
		result.Err = err
		result.ErrorReason = model.TransferErrorReasonInternal
		return result
	}

	// Transaction Process
	result = func() model.TransferResult {
		// lock accounts
		fromAccount, err := service.bankAccountLockService.Lock(ctx, req.FromBankAccountID, tx)
		if err != nil {
			result.Err = err
			result.ErrorReason = model.TransferErrorReasonInternal
			if model.IsNotFoundErrInInfraAdapter(err) {
				result.ErrorReason = model.TransferErrorReasonBankAccountNotFound
			}
			return result
		}
		toAccount, err := service.bankAccountLockService.Lock(ctx, req.ToBankAccountID, tx)
		if err != nil {
			result.Err = err
			result.ErrorReason = model.TransferErrorReasonInternal
			if model.IsNotFoundErrInInfraAdapter(err) {
				result.ErrorReason = model.TransferErrorReasonBankAccountNotFound
			}
			return result
		}

		// create transaction record
		record := model.NewTransferRecord(
			"uuid v4", fromAccount.ID, toAccount.ID, req.Money, time.Now(),
		)
		record.SetTransferStatus(model.TransferStatusInCheckBalance)
		err = service.transactionRecordRepository.Create(ctx, record.DTO())
		if err != nil {
			result.Err = err
			result.ErrorReason = model.TransferErrorReasonInternal
			return result
		}

		// check balance
		if fromAccount.Balance.IsLessThan(req.Money) {
			err = errors.New("insufficient balance")
			result.TransactionRecord.SetTransferStatus(model.TransferStatusInCheckBalanceFailed)
			err2 := service.transactionRecordRepository.Update(ctx, record.DTO())
			if err2 != nil {
				err = errors.Join(err, err2)
			}
			result.Err = err
			result.ErrorReason = model.TransferErrorReasonInsufficientBalance
			return result
		}

		// withdraw from the fromAccount
		result.TransactionRecord.SetTransferStatus(
			model.TransferStatusInTransfer,
		)
		result.TransactionRecord.UpdatedAt = time.Now()
		service.transactionRecordRepository.Update(ctx, result.TransactionRecord.DTO())
		fromAccount.Withdraw(req.Money, time.Now())
		err = service.bankAccountRepository.Update(fromAccount.DTO())
		if err != nil {
			result.Err = err
			result.ErrorReason = model.TransferErrorReasonInternal
			return result
		}

		// deposit to the toAccount
		toAccount.Deposit(req.Money, time.Now())
		err = service.bankAccountRepository.Update(toAccount.DTO())
		if err != nil {
			result.Err = err
			result.ErrorReason = model.TransferErrorReasonInternal
			return result
		}

		// set transfer completed
		result.TransactionRecord.SetTransferStatus(model.TransferStatusCompleted)
		result.TransactionRecord.UpdatedAt = time.Now()
		err = service.transactionRecordRepository.Update(ctx, result.TransactionRecord.DTO())
		if err != nil {
			result.Err = err
			result.ErrorReason = model.TransferErrorReasonInternal
			return result
		}

		// Transaction Commit
		err = service.transactionManager.Commit(ctx, tx)
		if err != nil {
			result.Err = err
			result.ErrorReason = model.TransferErrorReasonInternal
			return result
		}
		result.Err = nil
		result.ErrorReason = ""
		return result
	}()

	// If Failed, Rollback
	if result.Err != nil {
		tmpErr := service.transactionManager.Rollback(ctx, tx)
		if tmpErr != nil {
			result.Err = errors.Join(err, tmpErr)
		}
		return result
	}

	// success
	return result
}
