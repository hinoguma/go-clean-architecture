package appinfraadapterlayer

import (
	"app/threelayeredarchitecture/appinfralayer"
	"context"
)

type IsolationLevelDTO string

const (
	IsolationLevelDefault         IsolationLevelDTO = "Default"
	IsolationLevelReadUncommitted IsolationLevelDTO = "ReadUncommitted"
	IsolationLevelReadCommitted   IsolationLevelDTO = "ReadCommitted"
	IsolationLevelWriteCommitted  IsolationLevelDTO = "WriteCommitted"
	IsolationLevelRepeatableRead  IsolationLevelDTO = "RepeatableRead"
	IsolationLevelSnapshot        IsolationLevelDTO = "Snapshot"
	IsolationLevelSerializable    IsolationLevelDTO = "Serializable"
	IsolationLevelLinearizable    IsolationLevelDTO = "Linearizable"
)

func (value IsolationLevelDTO) IsolationLevel() appinfralayer.IsolationLevel {
	switch value {
	case IsolationLevelDefault:
		return appinfralayer.IsolationLevelDefault
	case IsolationLevelReadUncommitted:
		return appinfralayer.IsolationLevelReadUncommitted
	case IsolationLevelReadCommitted:
		return appinfralayer.IsolationLevelReadCommitted
	case IsolationLevelWriteCommitted:
		return appinfralayer.IsolationLevelWriteCommitted
	case IsolationLevelRepeatableRead:
		return appinfralayer.IsolationLevelRepeatableRead
	case IsolationLevelSnapshot:
		return appinfralayer.IsolationLevelSnapshot
	case IsolationLevelSerializable:
		return appinfralayer.IsolationLevelSerializable
	case IsolationLevelLinearizable:
		return appinfralayer.IsolationLevelLinearizable
	}
	return appinfralayer.IsolationLevelDefault
}

type BeginTransactionRequest struct {
	ID             *string
	IsolationLevel IsolationLevelDTO
	ReadOnly       bool
}

func NewBeginTransactionRequest() BeginTransactionRequest {
	return BeginTransactionRequest{
		ID:             nil,
		IsolationLevel: IsolationLevelDefault,
		ReadOnly:       false,
	}
}

type Transaction struct {
	ID string
}

func NewTransaction(id string) Transaction {
	return Transaction{ID: id}
}

type TransactionManagerAdapterIF interface {
	Begin(ctx context.Context, req BeginTransactionRequest) (Transaction, error)
	Rollback(ctx context.Context, tx Transaction) error
	Commit(ctx context.Context, tx Transaction) error
}

type transactionManagerAdapter struct {
	manager appinfralayer.TransactionManagerIF
}

func NewTransactionManagerAdapter(manager appinfralayer.TransactionManagerIF) TransactionManagerAdapterIF {
	return &transactionManagerAdapter{manager: manager}
}

func (adapter transactionManagerAdapter) Begin(ctx context.Context, req BeginTransactionRequest) (Transaction, error) {
	opts := appinfralayer.TxOptions{
		Isolation: req.IsolationLevel.IsolationLevel(),
		ReadOnly:  req.ReadOnly,
	}
	id := ""
	if req.ID == nil {
		// todo: uuid v4
	} else {
		id = *req.ID
	}
	tx, err := adapter.manager.Begin(ctx, id, &opts)
	if err != nil {
		return Transaction{}, err
	}
	return Transaction{ID: tx.GetID()}, err
}

func (adapter transactionManagerAdapter) Rollback(ctx context.Context, tx Transaction) error {
	return adapter.manager.Rollback(ctx, tx.ID)
}

func (adapter transactionManagerAdapter) Commit(ctx context.Context, tx Transaction) error {
	return adapter.manager.Commit(ctx, tx.ID)
}
