package usecase

import (
	"context"
	"errors"

	"app/threelayerdarchitecture/applogic/domain"
	"app/threelayerdarchitecture/appinfraadapter/repository"
)

type TransferRequest struct {
	InitiatingUserID string
	FromAccountID    string
	ToAccountID      string
	Amount           int64
	CorrelationID    string
}

type TransferResult struct {
	TransactionID string
}

var (
	ErrInvalidTransferRequest = errors.New("transfer use case: invalid request")
	ErrUnauthorizedUser       = errors.New("transfer use case: initiating user cannot act on source account")
)

type TransferUseCase struct {
	accounts   repository.AccountPort
	recorder   repository.TransactionPort
	transferer domain.TransferService
}

func NewTransferUseCase(
	accounts repository.AccountPort,
	recorder repository.TransactionPort,
	transferer domain.TransferService,
) TransferUseCase {
	return TransferUseCase{
		accounts:   accounts,
		recorder:   recorder,
		transferer: transferer,
	}
}

func (uc TransferUseCase) Execute(ctx context.Context, req TransferRequest) (TransferResult, error) {
	if err := validateRequest(req); err != nil {
		return TransferResult{}, err
	}

	fromAccountDTO, err := uc.accounts.Fetch(ctx, req.FromAccountID)
	if err != nil {
		return TransferResult{}, err
	}
	fromAccount, err := accountDTOToDomain(fromAccountDTO)
	if err != nil {
		return TransferResult{}, err
	}
	if fromAccount.UserID != req.InitiatingUserID {
		return TransferResult{}, ErrUnauthorizedUser
	}

	toAccountDTO, err := uc.accounts.Fetch(ctx, req.ToAccountID)
	if err != nil {
		return TransferResult{}, err
	}
	toAccount, err := accountDTOToDomain(toAccountDTO)
	if err != nil {
		return TransferResult{}, err
	}

	debited, credited, txn, err := uc.transferer.Transfer(fromAccount, toAccount, req.Amount)
	if err != nil {
		return TransferResult{}, err
	}

	if err := uc.accounts.Persist(ctx, domainToAccountDTO(debited)); err != nil {
		return TransferResult{}, err
	}
	if err := uc.accounts.Persist(ctx, domainToAccountDTO(credited)); err != nil {
		return TransferResult{}, err
	}

	txn.InitiatingUserID = req.InitiatingUserID
	txn.CorrelationID = req.CorrelationID

	transactionID, err := uc.recorder.Record(ctx, domainToTransactionDTO(txn))
	if err != nil {
		return TransferResult{}, err
	}

	return TransferResult{TransactionID: transactionID}, nil
}

func accountDTOToDomain(dto repository.AccountDTO) (domain.Account, error) {
	account := domain.Account{
		ID:       dto.ID,
		UserID:   dto.UserID,
		Balance:  dto.Balance,
		Currency: domain.Currency(dto.Currency),
	}
	if err := account.Validate(); err != nil {
		return domain.Account{}, err
	}
	return account, nil
}

func domainToAccountDTO(account domain.Account) repository.AccountDTO {
	return repository.AccountDTO{
		ID:       account.ID,
		UserID:   account.UserID,
		Balance:  account.Balance,
		Currency: string(account.Currency),
	}
}

func domainToTransactionDTO(txn domain.Transaction) repository.TransactionDTO {
	return repository.TransactionDTO{
		ID:               txn.ID,
		FromAccountID:    txn.FromAccountID,
		ToAccountID:      txn.ToAccountID,
		Amount:           txn.Amount,
		Currency:         string(txn.Currency),
		OccurredAtUnix:   txn.OccurredAt.Unix(),
		CorrelationID:    txn.CorrelationID,
		InitiatingUserID: txn.InitiatingUserID,
	}
}

func validateRequest(req TransferRequest) error {
	switch {
	case req.InitiatingUserID == "":
		return ErrInvalidTransferRequest
	case req.FromAccountID == "":
		return ErrInvalidTransferRequest
	case req.ToAccountID == "":
		return ErrInvalidTransferRequest
	case req.Amount <= 0:
		return ErrInvalidTransferRequest
	default:
		return nil
	}
}
