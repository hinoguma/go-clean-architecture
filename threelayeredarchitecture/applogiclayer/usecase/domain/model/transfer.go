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
	TransferStatusInCheckBalance TransferStatus = "IN_CHECK_BALANCE"
	TransferStatusInTransfer     TransferStatus = "IN_TRANSFER"
	TransferStatusCompleted      TransferStatus = "COMPLETED"

	TransferStatusInCheckBalanceFailed TransferStatus = "IN_CHECK_BALANCE_FAILED"
	TransferStatusInTransferFailed     TransferStatus = "IN_TRANSFER_FAILED"
	TransferStatusUnknownFailed        TransferStatus = "UNKNOWN_FAILED"
)

type TransferRecord struct{}

type TransactionRecord struct {
	ID                string
	Type              TransactionType
	FromBankAccountID string
	ToBankAccountID   string
	Money             Money
	TransferStatus    TransferStatus
	TransferLogs	   []TransferLog
	DataItem
}

func (record *TransactionRecord) SetTransferStatus(status TransferStatus) *TransactionRecord {
	record.TransferStatus = status
	return record
}

func (record *TransactionRecord) SetTransferCompleted() *TransactionRecord {
	record.TransferStatus = "COMPLETED"
	return record
}

func NewTransferRecord(
	id string,
	fromBankAccountID string,
	toBankAccountID string,
	money Money,
	createdAt time.Time,
) TransactionRecord {
	item := TransactionRecord{
		ID:                id,
		FromBankAccountID: fromBankAccountID,
		ToBankAccountID:   toBankAccountID,
		Money:             money,
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

type TransferLog struct{
	ActionType TransferActionType
	ActionStatus TransferActionStatus
	Timestamp time.Time
}



type TransferRequest struct {
	FromBankAccountID string
	ToBankAccountID   string
	Money             Money
}

type TransferResult struct {
	TransactionRecord TransactionRecord
}
