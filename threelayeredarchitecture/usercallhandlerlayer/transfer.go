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
	// todo: how to handle error cases?

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
	adapterRes, err := handler.adapter.Execute(ctx, adapterReq)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{}, err
	}

	// convert adapter response to API Gateway response format
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Body:       fmt.Sprint(`{"transactionId":%s}`, adapterRes.TransactionID),
	}, nil
}

type transferRequstBodyJson struct {
	FromBankAccountID string `json:"fromBankAccountId"`
	ToBankAccountID   string `json:"toBankAccountId"`
	Amount            int    `json:"amount"`
	Currency          string `json:"currency"`
}
