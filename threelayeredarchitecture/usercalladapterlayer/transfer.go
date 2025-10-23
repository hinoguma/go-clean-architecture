package usercalladapterlayer

import (
	"app/threelayeredarchitecture/applogiclayer/usecase"
	"app/threelayeredarchitecture/applogiclayer/usecase/domain/model"
	"context"
)

type TransferAdapterRequest struct {
	FromBankAccountID string
	ToBankAccountID   string
	Amount            int64
	Currency          string
}

type TransferAdapterResponse struct {
	TransactionID string
	ErrorReason   TransferAdapterErrorReason
	Err           error
}

func NewTransferAdapterResponseByUsecaseRes(ucRes usecase.TransferUsecaseResponse) TransferAdapterResponse {
	if ucRes.Err != nil {
		// success
		return TransferAdapterResponse{
			TransactionID: ucRes.TransactionRecord.ID,
		}
	}
	// error
	return TransferAdapterResponse{
		TransactionID: "",
		Err:           ucRes.Err,
		ErrorReason: TransferAdapterErrorReason{
			AuthError:           ucRes.ErrReason.AuthenticateError,
			ValidateError:       ucRes.ErrReason.ValidateError,
			BankAccountNotFound: ucRes.ErrReason.BankAccountNotFound,
			InsufficientBalance: ucRes.ErrReason.InsufficientBalance,
			InternalError:       ucRes.ErrReason.InternalError,
		},
	}
}

type TransferAdapterErrorReason struct {
	AuthError           bool
	ValidateError       bool
	BankAccountNotFound bool
	InsufficientBalance bool
	InternalError       bool
}

type TransferAdapterIF interface {
	Execute(ctx context.Context, req TransferAdapterRequest) TransferAdapterResponse
}

type TransferAdapter struct {
	// usecase
	usecase usecase.TransferUsecaseIF
}

func NewTransferAdapter(usecase usecase.TransferUsecaseIF) *TransferAdapter {
	return &TransferAdapter{
		usecase: usecase,
	}
}

func (th *TransferAdapter) Execute(ctx context.Context, req TransferAdapterRequest) TransferAdapterResponse {

	// convert adapter request to usecase request format
	usecaseReq := usecase.TransferUsecaseRequest{
		FromBankAccountID: req.FromBankAccountID,
		ToBankAccountID:   req.ToBankAccountID,
		Money: model.Money{
			Amount:   int64(req.Amount),
			Currency: model.Currency(req.Currency),
		},
	}

	// call adapter logic
	usecaseRes := th.usecase.Execute(ctx, usecaseReq)
	return NewTransferAdapterResponseByUsecaseRes(usecaseRes)
}
