package usercallhandlerlayer

import (
	"app/threelayeredarchitecture/usercalladapterlayer"
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

type TransferHandlerIF interface {
	Execute(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error)
}

type TransferHandler struct {
	// adapter
	adapter usercalladapterlayer.TransferAdapterIF
}

func NewTransferHandler(adapter usercalladapterlayer.TransferAdapterIF) *TransferHandler {
	return &TransferHandler{
		adapter: adapter,
	}
}

func (handler *TransferHandler) Execute(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	// convert API Gateway request to adapter request format
	var reqBody transferRequstBodyJson
	err := json.Unmarshal([]byte(event.Body), &reqBody)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{}, err
	}
	adapterReq := usercalladapterlayer.TransferAdapterRequest{
		FromBankAccountID: reqBody.FromBankAccountID,
		ToBankAccountID:   reqBody.ToBankAccountID,
		Amount:            reqBody.Amount,
		Currency:          reqBody.Currency,
	}

	// call adapter logic
	adapterRes := handler.adapter.Execute(ctx, adapterReq)

	// response
	return handler.responseByAdapterResponse(adapterRes)
}

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
		// todo: how to get validation error details
		return NewValidateErrResponse(adapterRes.Err), adapterRes.Err
	}
	if adapterRes.ErrorReason.BankAccountNotFound {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 460,
		}, adapterRes.Err
	}
	if adapterRes.ErrorReason.InsufficientBalance {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 461,
		}, adapterRes.Err
	}
	// internal error
	return NewInternalServerErrResponse(adapterRes.Err), adapterRes.Err
}

type transferRequstBodyJson struct {
	FromBankAccountID string `json:"fromBankAccountId"`
	ToBankAccountID   string `json:"toBankAccountId"`
	Amount            int    `json:"amount"`
	Currency          string `json:"currency"`
}
