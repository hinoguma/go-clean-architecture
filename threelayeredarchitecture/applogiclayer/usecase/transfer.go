package usecase

import (
	"app/threelayeredarchitecture/applogiclayer/usecase/domain/model"
	"app/threelayeredarchitecture/applogiclayer/usecase/domain/service"
	"app/threelayeredarchitecture/crosscuttinglayer"
	"context"
)

type TransferUsecaseRequest struct {
	IdempotencyKey    string
	AuthToken         string
	FromBankAccountID string
	ToBankAccountID   string
	Money             model.Money
}

type TransferUsecaseResponse struct {
	TransactionRecord model.TransactionRecord
	ErrReason         TransferUsecaseErrorReason
	Err               error
}

type TransferUsecaseErrorReason struct {
	AuthenticateError   bool
	ValidateError       bool
	BankAccountNotFound bool
	InsufficientBalance bool
	InternalError       bool
}

type TransferUsecaseIF interface {
	Execute(ctx context.Context, req TransferUsecaseRequest) TransferUsecaseResponse
}

type TransferUsecase struct {
	authService     service.BankCustomerAuthServiceIF
	transferService service.TransferServiceIF
}

func NewTransferUsecase(
	authService service.BankCustomerAuthServiceIF,
	transferService service.TransferServiceIF,
) TransferUsecaseIF {
	return &TransferUsecase{
		authService:     authService,
		transferService: transferService,
	}
}

func (uc *TransferUsecase) Execute(ctx context.Context, req TransferUsecaseRequest) TransferUsecaseResponse {
	// Authentication
	nowTs := crosscuttinglayer.NowTs()
	authRes := uc.authService.Authenticate(ctx, model.BankCustomerAuthRequest{Token: req.AuthToken, Timestamp: nowTs.Int64()})
	if !authRes.IsSuccess() {
		reason := TransferUsecaseErrorReason{}
		if authRes.IsAuthenticateFailedError() {
			reason.AuthenticateError = true
		} else {
			reason.InternalError = true
		}
		return TransferUsecaseResponse{
			ErrReason: reason,
			Err:       authRes.Err,
		}
	}

	// validate request
	appLogicErr := uc.validateRequest(req)
	if appLogicErr != nil {
		return TransferUsecaseResponse{
			ErrReason: TransferUsecaseErrorReason{ValidateError: true},
			Err:       appLogicErr,
		}
	}

	// Transfer
	transferReq := model.TransferRequest{
		IdempotencyKey:    req.IdempotencyKey,
		FromBankAccountID: req.FromBankAccountID,
		ToBankAccountID:   req.ToBankAccountID,
		Money:             req.Money,
	}
	transferRes := uc.transferService.Execute(ctx, transferReq)

	// Error Handling
	if transferRes.Err != nil {
		reason := TransferUsecaseErrorReason{}
		switch transferRes.ErrorReason {
		case model.TransferErrorReasonBankAccountNotFound:
			reason.BankAccountNotFound = true
		case model.TransferErrorReasonInsufficientBalance:
			reason.InsufficientBalance = true
		default:
			reason.InternalError = true
		}
		return TransferUsecaseResponse{
			ErrReason: reason,
			Err:       transferRes.Err,
		}
	}

	// Success
	return TransferUsecaseResponse{
		TransactionRecord: transferRes.TransactionRecord,
		Err:               nil,
	}
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
