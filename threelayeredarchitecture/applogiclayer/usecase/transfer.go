package usecase

import (
	"app/threelayeredarchitecture/applogiclayer/usecase/domain/model"
	"app/threelayeredarchitecture/applogiclayer/usecase/domain/service"
	"context"
)

type TransferUsecaseRequest struct {
	FromBankAccountID string
	ToBankAccountID   string
	Money             model.Money
}

type TransferUsecaseResponse struct {
	TransactionRecord model.TransactionRecord
}

type TransferUsecaseIF interface {
	Execute(ctx context.Context, req TransferUsecaseRequest) (TransferUsecaseResponse, error)
}

type TransferUsecase struct {
	transferService service.TransferServiceIF
}

func NewTransferUsecase(
	transferService service.TransferServiceIF,
) TransferUsecaseIF {
	return &TransferUsecase{
		transferService: transferService,
	}
}

func (uc *TransferUsecase) Execute(ctx context.Context, req TransferUsecaseRequest) (TransferUsecaseResponse, error) {

	// validate request
	appLogicErr := uc.validateRequest(req)
	if appLogicErr != nil {
		return TransferUsecaseResponse{}, appLogicErr
	}

	transferReq := model.TransferRequest{
		FromBankAccountID: req.FromBankAccountID,
		ToBankAccountID:   req.ToBankAccountID,
		Money:             req.Money,
	}
	transferRes, err := uc.transferService.Execute(transferReq)
	if err != nil {
		return TransferUsecaseResponse{}, err
	}

	// error patterns
	// - validation error
	// - bank account not found
	// - bank customer not found
	// - insufficient balance
	// - internal error

	return TransferUsecaseResponse{
		TransactionRecord: transferRes.TransactionRecord,
	}, nil
}

func (uc TransferUsecase) validateRequest(req TransferUsecaseRequest) *model.AppLogicError {
	details := make([]model.ValidationErrorDetail, 0)
	if req.FromBankAccountID == "" {
		details = append(details, model.NewRequiredErrDetail("fromBankAccountId"))
	}
	if req.ToBankAccountID == "" {
		details = append(details, model.NewRequiredErrDetail("toBankAccountId"))
	}
	if !req.Money.IsValidAmount() {
		details = append(details, model.NewMinValueErrDetail("amount", 1))
	}
	if !req.Money.IsValidCurrency() {
		details = append(details, model.NewRequiredErrDetail("currency"))
	}

	if len(details) > 0 {
		return model.NewValidationError(details)
	}
	return nil
}
