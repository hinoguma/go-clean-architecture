package usercallhandlerlayer

import (
	"app/threelayeredarchitecture/usercalladapterlayer"
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

/*
**************************************

		Request
	 **************************************
*/
type transferRequstBodyJson struct {
	FromBankAccountID *string   `json:"fromBankAccountId,omitempty"`
	ToBankAccountID   *string   `json:"toBankAccountId,omitempty"`
	Amount            *int64    `json:"amount,omitempty"`
	Currency          *Currency `json:"currency,omitempty"`
}

func (body transferRequstBodyJson) Validate() []ValidationErrorDetail {
	details := make([]ValidationErrorDetail, 0)

	if body.FromBankAccountID == nil {
		details = append(details, NewRequiredValidationError("fromBankAccountId"))
	} else if *body.FromBankAccountID == "" {
		details = append(details, NewEmptyError("fromBankAccountId"))
	}

	if body.ToBankAccountID == nil {
		details = append(details, NewRequiredValidationError("toBankAccountId"))
	} else if *body.ToBankAccountID == "" {
		details = append(details, NewEmptyError("toBankAccountId"))
	}

	if body.Amount == nil {
		details = append(details, NewRequiredValidationError("amount"))
	} else if *body.Amount <= 0 || *body.Amount > 50*10000 {
		details = append(details, NewMustBeRangeError("amount", 1, 50*10000))
	}

	if body.Currency == nil {
		details = append(details, NewRequiredValidationError("currency"))
	} else if !body.Currency.IsValid() {
		details = append(details, NewInvalidValueError("currency"))
	} else {

	}

	return details
}

func (body transferRequstBodyJson) AdapterRequest() usercalladapterlayer.TransferAdapterRequest {
	return usercalladapterlayer.TransferAdapterRequest{
		FromBankAccountID: body.GetFromBankAccountID(),
		ToBankAccountID:   body.GetToBankAccountID(),
		Amount:            body.GetAmount(),
		Currency:          string(body.GetCurrency()),
	}
}

func (body transferRequstBodyJson) GetFromBankAccountID() string {
	if body.FromBankAccountID == nil {
		return ""
	}
	return *body.FromBankAccountID
}

func (body transferRequstBodyJson) GetToBankAccountID() string {
	if body.ToBankAccountID == nil {
		return ""
	}
	return *body.ToBankAccountID
}

func (body transferRequstBodyJson) GetAmount() int64 {
	if body.Amount == nil {
		return 0
	}
	return *body.Amount
}

func (body transferRequstBodyJson) GetCurrency() Currency {
	if body.Currency == nil {
		return ""
	}
	return *body.Currency
}

/***************************************
	Handler
 ***************************************/

type TransferHandlerIF interface {
	Execute(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error)
}

type TransferHandler struct {
	adapter usercalladapterlayer.TransferAdapterIF
}

func NewTransferHandler(adapter usercalladapterlayer.TransferAdapterIF) *TransferHandler {
	return &TransferHandler{
		adapter: adapter,
	}
}

func (handler *TransferHandler) Execute(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	// convert raw body to body struct
	var reqBody transferRequstBodyJson
	err := json.Unmarshal([]byte(event.Body), &reqBody)
	if err != nil {
		return NewCannotMarshalJsonErrResponse(), err
	}

	// validate request
	errDetails := reqBody.Validate()
	if len(errDetails) > 0 {
		return NewValidateErrResponseWithDetail(errDetails), nil
	}

	// convert to adapter request
	adapterReq := reqBody.AdapterRequest()

	// call adapter logic
	adapterRes := handler.adapter.Execute(ctx, adapterReq)

	// response
	return handler.responseByAdapterResponse(adapterRes)
}

/***************************************
	Response
 ***************************************/

const (
	CodeBankAccountNotFound = 460
	CodeInsufficientBalance = 461
)

func (handler TransferHandler) responseByAdapterResponse(adapterRes usercalladapterlayer.TransferAdapterResponse) (events.APIGatewayV2HTTPResponse, error) {
	if adapterRes.Err == nil {
		// success
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 200,
			Body:       fmt.Sprint(`{"transactionId":%s}`, adapterRes.TransactionID),
		}, nil
	}

	if adapterRes.ErrorReason.AuthError {
		return NewAuthErrResponse(adapterRes.Err), adapterRes.Err
	}

	if adapterRes.ErrorReason.ValidateError {
		return NewValidateErrResponseUnderAdapter(), adapterRes.Err
	}

	if adapterRes.ErrorReason.BankAccountNotFound {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: CodeBankAccountNotFound,
		}, adapterRes.Err
	}

	if adapterRes.ErrorReason.InsufficientBalance {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: CodeInsufficientBalance,
		}, adapterRes.Err
	}
	// internal error
	return NewInternalServerErrResponse(adapterRes.Err), adapterRes.Err
}
