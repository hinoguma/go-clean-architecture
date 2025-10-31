package model

import (
	"app/threelayeredarchitecture/appinfraadapterlayer"
	"time"
)

type TransactionType string

const (
	TransactionTypeTransfer TransactionType = "TRANSFER"
)

type TransferStatus string

const (
	TransferStatusPrepare         TransferStatus = "PREPARE"
	TransferStatusFromBankAccount TransferStatus = "LOCK_FROM_BANK_ACCOUNT"
	TransferStatusToBankAccount   TransferStatus = "LOCK_TO_BANK_ACCOUNT"
	TransferStatusInCheckBalance  TransferStatus = "IN_CHECK_BALANCE"
	TransferStatusInTransfer      TransferStatus = "IN_TRANSFER"
	TransferStatusCompleted       TransferStatus = "COMPLETED"

	TransferStatusFromBankAccountFailed TransferStatus = "FROM_BANK_ACCOUNT_FAILED"
	TransferStatusToBankAccountFailed   TransferStatus = "TO_BANK_ACCOUNT_FAILED"
	TransferStatusInCheckBalanceFailed  TransferStatus = "IN_CHECK_BALANCE_FAILED"
	TransferStatusInTransferFailed      TransferStatus = "IN_TRANSFER_FAILED"
	TransferStatusUnknownFailed         TransferStatus = "UNKNOWN_FAILED"
)

type TransferRecord struct{}

type TransactionRecord struct {
	ID                string
	IdempotencyKey    string
	Type              TransactionType
	FromBankAccountID string
	ToBankAccountID   string
	Money             Money
	TransferStatus    TransferStatus
	DataItem
}

func (record *TransactionRecord) SetFromDTO(dto appinfraadapterlayer.TransactionRecordDTO) *TransactionRecord {
	record.ID = dto.ID
	record.IdempotencyKey = dto.IdempotencyKey
	record.Type = TransactionType(dto.Type)
	record.FromBankAccountID = dto.FromBankAccountID
	record.ToBankAccountID = dto.ToBankAccountID
	record.Money.Amount = dto.Amount
	record.Money.Currency = Currency(dto.Currency)
	record.TransferStatus = TransferStatus(dto.TransferStatus)
	record.DataItem.SetFromDTO(dto.DatabaseItem)
	return record
}

func (record *TransactionRecord) SetTransferStatus(status TransferStatus) *TransactionRecord {
	record.TransferStatus = status
	return record
}

func (record *TransactionRecord) SetTransferCompleted() *TransactionRecord {
	record.TransferStatus = TransferStatusCompleted
	return record
}

func (record TransactionRecord) IsCompleted() bool {
	return record.TransferStatus == TransferStatusCompleted
}

func (record TransactionRecord) IsConsistentRequest(req TransferRequest) bool {
	return record.IdempotencyKey == req.IdempotencyKey &&
		record.FromBankAccountID == req.FromBankAccountID &&
		record.ToBankAccountID == req.ToBankAccountID &&
		record.Money.IsEqual(req.Money)
}

func NewTransferRecord(
	id string,
	req TransferRequest,
	createdAt time.Time,
) TransactionRecord {
	item := TransactionRecord{
		ID:                id,
		TransferStatus:    TransferStatusPrepare,
		IdempotencyKey:    req.IdempotencyKey,
		FromBankAccountID: req.FromBankAccountID,
		ToBankAccountID:   req.ToBankAccountID,
		Money:             req.Money,
	}
	item.CreatedAt = createdAt
	return item
}

func (record TransactionRecord) DTO() appinfraadapterlayer.TransactionRecordDTO {
	dto := appinfraadapterlayer.TransactionRecordDTO{
		ID:                record.ID,
		FromBankAccountID: record.FromBankAccountID,
		ToBankAccountID:   record.ToBankAccountID,
		Amount:            record.Money.Amount,
		Currency:          string(record.Money.Currency),
	}
	dto.DatabaseItem.CreatedAt = record.CreatedAt.Unix()
	dto.DatabaseItem.UpdatedAt = record.UpdatedAt.Unix()
	return dto
}

type TransferActionType string

const (
	TransferActionTypeCheckBalance TransferActionType = "CHECK_BALANCE"
	TransferActionTypeWithdraw     TransferActionType = "WITHDRAW"
	TransferActionTypeDeposit      TransferActionType = "DEPOSIT"
)

type TransferActionStatus string

const (
	TransferActionStatusStarted    TransferActionStatus = "STARTED"
	TransferActionStatusInProgress TransferActionStatus = "IN_PROGRESS"
	TransferActionStatusCompleted  TransferActionStatus = "COMPLETED"
	TransferActionStatusFailed     TransferActionStatus = "FAILED"
)

type TransferRequest struct {
	IdempotencyKey    string
	FromBankAccountID string
	ToBankAccountID   string
	Money             Money
}

type TransferResult struct {
	TransactionRecord TransactionRecord
	ErrorReason       TransferErrorReason
	Err               error
}

type TransferErrorReason string

func (value TransferErrorReason) TransferStatus() TransferStatus {
	switch value {
	case TransferErrorReasonNotConsistRequest:
		return TransferStatusUnknownFailed
	case TransferErrorReasonFromBankAccountNotFound:
		return TransferStatusFromBankAccountFailed
	case TransferErrorReasonToBankAccountNotFound:
		return TransferStatusToBankAccountFailed
	case TransferErrorReasonBankCustomerNotFound:
		return TransferStatusUnknownFailed
	case TransferErrorReasonInsufficientBalance:
		return TransferStatusInCheckBalanceFailed
	case TransferErrorReasonInternal:
		return TransferStatusUnknownFailed
	}
	return TransferStatusUnknownFailed
}

const (
	TransferErrorReasonNotConsistRequest       TransferErrorReason = "NOT_CONSIST_REQUEST"
	TransferErrorReasonFromBankAccountNotFound TransferErrorReason = "FROM_BANK_ACCOUNT_NOT_FOUND"
	TransferErrorReasonToBankAccountNotFound   TransferErrorReason = "TO_BANK_ACCOUNT_NOT_FOUND"
	TransferErrorReasonBankCustomerNotFound    TransferErrorReason = "BANK_CUSTOMER_NOT_FOUND"
	TransferErrorReasonInsufficientBalance     TransferErrorReason = "INSUFFICIENT_BALANCE"
	TransferErrorReasonInternal                TransferErrorReason = "INTERNAL_ERROR"
)
