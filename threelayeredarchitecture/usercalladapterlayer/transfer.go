package usercalladapterlayer

import (
	"app/threelayeredarchitecture/applogiclayer/usecase"
	"app/threelayeredarchitecture/applogiclayer/usecase/domain/model"
	"context"
)

type TransferAdapterRequest struct {
	FromBankAccountID string
	ToBankAccountID   string
	Amount            int
	Currency          string
}

type TransferAdapterResponse struct {
	TransactionID string
}

type TransferAdapterIF interface {
	Execute(ctx context.Context, req TransferAdapterRequest) (TransferAdapterResponse, error)
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

func (th *TransferAdapter) Execute(ctx context.Context, req TransferAdapterRequest) (TransferAdapterResponse, error) {

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
	usecaseRes, err := th.usecase.Execute(usecaseReq)
	if err != nil {
		return TransferAdapterResponse{}, err
	}

	// convert usecase response to adapter response format
	return TransferAdapterResponse{
		TransactionID: usecaseRes.TransactionID,
	}, nil
}
