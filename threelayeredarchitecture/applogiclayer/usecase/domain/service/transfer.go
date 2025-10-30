package service

import (
	"app/threelayeredarchitecture/appinfraadapterlayer"
	"app/threelayeredarchitecture/applogiclayer/usecase/domain/model"
	"app/threelayeredarchitecture/crosscuttinglayer"
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

	// check for idempotency
	recordDTO, err := service.transactionRecordRepository.GetByIdempotencyKey(ctx, req.IdempotencyKey)
	isNotFoundIdempotency := crosscuttinglayer.IsDataNotFound(err)
	if err != nil && !isNotFoundIdempotency {
		result.Err = crosscuttinglayer.ErrLift(err, ctx)
		result.ErrorReason = model.TransferErrorReasonInternal
		return result
	}
	// if not found, make transaction record
	if !isNotFoundIdempotency {
		result.TransactionRecord.SetFromDTO(recordDTO)
		if !result.TransactionRecord.IsConsistantRequest(req) {
			appErr := crosscuttinglayer.NewAppUtilError("not consist request", ctx)
			appErr.Attr("record", result.TransactionRecord).
				Attr("request", req)
			result.Err = appErr
			result.ErrorReason = model.TransferErrorReasonNotConsistRequest
			return result
		}
		if result.TransactionRecord.IsCompleted() {
			result.Err = nil
			result.ErrorReason = ""
			return result
		}
	}
	// create transaction record
	result.TransactionRecord = model.NewTransferRecord(
		crosscuttinglayer.IssueRandomStrID(), req, time.Now(),
	)
	err = service.transactionRecordRepository.Create(ctx, result.TransactionRecord.DTO())
	if err != nil {
		result.Err = crosscuttinglayer.ErrLift(err, ctx)
		result.ErrorReason = model.TransferErrorReasonInternal
		return result
	}
	// if found
	// if request information does not match, return error

	// Transaction Begins
	txBeginReq := appinfraadapterlayer.NewBeginTransactionRequest()
	tx, err := service.transactionManager.Begin(ctx, txBeginReq)
	if err != nil {
		result.Err = crosscuttinglayer.ErrLift(err, ctx)
		result.ErrorReason = model.TransferErrorReasonInternal
		return result
	}

	// Transaction Process
	result = func() model.TransferResult {
		// lock accounts
		fromAccount, err := service.bankAccountLockService.Lock(ctx, req.FromBankAccountID, tx)
		if err != nil {
			result.Err = crosscuttinglayer.ErrLift(err, ctx)
			result.ErrorReason = model.TransferErrorReasonInternal
			if model.IsNotFoundErrInInfraAdapter(err) {
				result.ErrorReason = model.TransferErrorReasonBankAccountNotFound
			}
			return result
		}
		toAccount, err := service.bankAccountLockService.Lock(ctx, req.ToBankAccountID, tx)
		if err != nil {
			result.Err = crosscuttinglayer.ErrLift(err, ctx)
			result.ErrorReason = model.TransferErrorReasonInternal
			if model.IsNotFoundErrInInfraAdapter(err) {
				result.ErrorReason = model.TransferErrorReasonBankAccountNotFound
			}
			return result
		}

		// check balance
		if fromAccount.Balance.IsLessThan(req.Money) {
			err = errors.New("insufficient balance")
			result.TransactionRecord.SetTransferStatus(model.TransferStatusInCheckBalanceFailed)
			err2 := service.transactionRecordRepository.TxPut(ctx, record.DTO(), tx)
			if err2 != nil {
				err = errors.Join(err, err2)
			}
			result.Err = crosscuttinglayer.ErrLift(err, ctx)
			result.ErrorReason = model.TransferErrorReasonInsufficientBalance
			return result
		}

		// withdraw from the fromAccount
		result.TransactionRecord.SetTransferStatus(
			model.TransferStatusInTransfer,
		)
		result.TransactionRecord.UpdatedAt = time.Now()
		err = service.transactionRecordRepository.TxPut(ctx, result.TransactionRecord.DTO(), tx)
		if err != nil {
			result.Err = crosscuttinglayer.ErrLift(err, ctx)
			result.ErrorReason = model.TransferErrorReasonInternal
			return result
		}
		fromAccount.Withdraw(req.Money, time.Now())
		err = service.bankAccountRepository.TxPut(ctx, fromAccount.DTO(), tx)
		if err != nil {
			result.Err = crosscuttinglayer.ErrLift(err, ctx)
			result.ErrorReason = model.TransferErrorReasonInternal
			return result
		}

		// deposit to the toAccount
		toAccount.Deposit(req.Money, time.Now())
		err = service.bankAccountRepository.TxPut(ctx, toAccount.DTO(), tx)
		if err != nil {
			result.Err = crosscuttinglayer.ErrLift(err, ctx)
			result.ErrorReason = model.TransferErrorReasonInternal
			return result
		}

		// set transfer completed
		result.TransactionRecord.SetTransferStatus(model.TransferStatusCompleted)
		result.TransactionRecord.UpdatedAt = time.Now()
		err = service.transactionRecordRepository.TxPut(ctx, result.TransactionRecord.DTO(), tx)
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
