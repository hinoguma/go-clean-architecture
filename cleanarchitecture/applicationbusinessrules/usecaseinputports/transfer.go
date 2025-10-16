package usecaseinputports

import "context"

type TransferInput struct {
	BankAccountUserID string
	FromBankAccountID string
	ToBankAccountID   string
	Amount            int64
}

type TransferUseCaseInputPort interface {
	Execute(ctx context.Context, input TransferInput) error
}
