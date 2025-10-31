package service

import (
	"app/threelayeredarchitecture/appinfraadapterlayer"
	"app/threelayeredarchitecture/applogiclayer/usecase/domain/model"
	"app/threelayeredarchitecture/crosscuttinglayer"
	"context"
	"errors"
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
	// make sure idempotency
	result, hasDone := service.makeSureIdempotency(ctx, req)
	if hasDone {
		return result
	}
	if result.Err != nil {
		return result
	}

	// transaction begins and lock transaction record
	txPtr, record, err := service.beginTransaction(ctx, result.TransactionRecord)
	if err != nil {
		if txPtr != nil {
			tmpErr := service.transactionManager.Rollback(ctx, *txPtr)
			if tmpErr != nil {
				err = errors.Join(err, tmpErr)
			}
		}
		result.Err = err
		result.ErrorReason = model.TransferErrorReasonInternal
		return result
	}

	// transaction process
	tx := *txPtr
	result = service.executeTransactionProcess(ctx, record, tx)
	if result.Err == nil {
		// success
		return result
	}

	// if failed, rollback
	tmpErr := service.transactionManager.Rollback(ctx, tx)
	if tmpErr != nil {
		result.Err = errors.Join(err, tmpErr)
	}
	// update transaction record status to failed
	result.TransactionRecord.UpdatedAt = crosscuttinglayer.Now()
	result.TransactionRecord.SetTransferStatus(
		result.ErrorReason.TransferStatus(),
	)
	tmpErr = service.transactionRecordRepository.Put(ctx, result.TransactionRecord.DTO())
	if tmpErr != nil {
		result.Err = errors.Join(err, tmpErr)
	}
	// error
	return result
}

func (service *TransferService) makeSureIdempotency(
	ctx context.Context, req model.TransferRequest,
) (model.TransferResult, bool) {
	result := model.TransferResult{}
	// check for idempotency
	recordDTO, err := service.transactionRecordRepository.GetByIdempotencyKey(ctx, req.IdempotencyKey)
	isNotFoundIdempotency := !crosscuttinglayer.IsDataNotFound(err)
	if err != nil && !isNotFoundIdempotency {
		result.Err = crosscuttinglayer.ErrLift(err, ctx)
		result.ErrorReason = model.TransferErrorReasonInternal
		return result, false
	}
	// if not found, make transaction record
	if !isNotFoundIdempotency {
		// found idempotency
		record := model.TransactionRecord{}
		record.SetFromDTO(recordDTO)
		if !record.IsConsistentRequest(req) {
			appErr := crosscuttinglayer.NewAppUtilError("not consist request", ctx)
			appErr.Attr("record", result.TransactionRecord).
				Attr("request", req)
			result.TransactionRecord = record
			result.Err = appErr
			result.ErrorReason = model.TransferErrorReasonNotConsistRequest
			return result, false
		}
		// return the same result if already completed
		if result.TransactionRecord.IsCompleted() {
			result.TransactionRecord = record
			result.Err = nil
			result.ErrorReason = ""
			return result, true
		}
		// if not completed, maybe failed before, return the record and go to retry
		return model.TransferResult{
			TransactionRecord: record,
			Err:               nil,
			ErrorReason:       "",
		}, false
	}

	// not found idempotency -> first time request
	// create transaction record
	record := model.NewTransferRecord(
		crosscuttinglayer.IssueRandomStrID(), req, crosscuttinglayer.Now(),
	)
	err = service.transactionRecordRepository.Create(ctx, record.DTO())
	if err != nil {
		result.Err = crosscuttinglayer.ErrLift(err, ctx)
		result.ErrorReason = model.TransferErrorReasonInternal
		return result, false
	}

	return model.TransferResult{
		TransactionRecord: record,
		Err:               nil,
		ErrorReason:       "",
	}, false
}

func (service *TransferService) beginTransaction(
	ctx context.Context, record model.TransactionRecord,
) (*appinfraadapterlayer.Transaction, model.TransactionRecord, error) {
	txBeginReq := appinfraadapterlayer.NewBeginTransactionRequest()
	tx, err := service.transactionManager.Begin(ctx, txBeginReq)
	if err != nil {
		return nil, model.TransactionRecord{}, crosscuttinglayer.ErrLift(err, ctx)
	}
	// lock transaction record
	recordDTO, err := service.transactionRecordRepository.Lock(ctx, record.ID, tx)
	if err != nil {
		return &tx, model.TransactionRecord{}, crosscuttinglayer.ErrLift(err, ctx)
	}
	record.SetFromDTO(recordDTO)
	return &tx, record, nil
}

func (service *TransferService) executeTransactionProcess(
	ctx context.Context, record model.TransactionRecord, tx appinfraadapterlayer.Transaction,
) model.TransferResult {
	result := model.TransferResult{
		TransactionRecord: record,
		ErrorReason:       model.TransferErrorReasonInternal,
		Err:               errors.New("unknown error"),
	}

	// lock accounts
	fromAccount, err := service.bankAccountLockService.Lock(ctx, record.FromBankAccountID, tx)
	if err != nil {
		result.Err = crosscuttinglayer.ErrLift(err, ctx)
		if model.IsNotFoundErrInInfraAdapter(err) {
			result.ErrorReason = model.TransferErrorReasonFromBankAccountNotFound
		}
		return result
	}
	toAccount, err := service.bankAccountLockService.Lock(ctx, record.ToBankAccountID, tx)
	if err != nil {
		result.Err = crosscuttinglayer.ErrLift(err, ctx)
		if model.IsNotFoundErrInInfraAdapter(err) {
			result.ErrorReason = model.TransferErrorReasonToBankAccountNotFound
		}
		return result
	}

	// check balance
	if fromAccount.Balance.IsLessThan(record.Money) {
		err = errors.New("insufficient balance")
		result.Err = crosscuttinglayer.ErrLift(err, ctx)
		result.ErrorReason = model.TransferErrorReasonInsufficientBalance
		return result
	}

	// withdraw from the fromAccount
	fromAccount.Withdraw(record.Money, crosscuttinglayer.Now())
	err = service.bankAccountRepository.TxPut(ctx, fromAccount.DTO(), tx)
	if err != nil {
		result.Err = crosscuttinglayer.ErrLift(err, ctx)
		result.ErrorReason = model.TransferErrorReasonInternal
		return result
	}

	// deposit to the toAccount
	toAccount.Deposit(record.Money, crosscuttinglayer.Now())
	err = service.bankAccountRepository.TxPut(ctx, toAccount.DTO(), tx)
	if err != nil {
		result.Err = crosscuttinglayer.ErrLift(err, ctx)
		result.ErrorReason = model.TransferErrorReasonInternal
		return result
	}

	// set transfer completed
	record.SetTransferStatus(model.TransferStatusCompleted)
	record.UpdatedAt = crosscuttinglayer.Now()
	err = service.transactionRecordRepository.TxPut(ctx, record.DTO(), tx)
	if err != nil {
		result.Err = crosscuttinglayer.ErrLift(err, ctx)
		result.ErrorReason = model.TransferErrorReasonInternal
		return result
	}

	// Transaction Commit
	err = service.transactionManager.Commit(ctx, tx)
	if err != nil {
		result.Err = crosscuttinglayer.ErrLift(err, ctx)
		result.ErrorReason = model.TransferErrorReasonInternal
		return result
	}

	// success
	result.TransactionRecord = record
	result.Err = nil
	result.ErrorReason = ""
	return result
}
