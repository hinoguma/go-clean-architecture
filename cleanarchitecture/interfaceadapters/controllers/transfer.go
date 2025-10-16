package controllers

import (
	"app/cleanarchitecture/applicationbusinessrules/usecaseinputports"
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
)

type TransferController struct {
	usecase usecaseinputports.TransferUseCaseInputPort
}

func NewTransferController(
	usecase usecaseinputports.TransferUseCaseInputPort,
) TransferController {
	return TransferController{
		usecase: usecase,
	}
}

func (controller TransferController) Execute(ctx context.Context, input events.APIGatewayV2HTTPRequest) error {
	// Convert http.Request to usecaseinputports.TransferInput
	bodyBytes := []byte(input.Body)
	var requestData struct {
		BankAccountUserID string `json:"bank_account_user_id"`
		FromBankAccountID string `json:"from_bank_account_id"`
		ToBankAccountID   string `json:"to_bank_account_id"`
		Amount            int64  `json:"amount"`
	}
	err := json.Unmarshal(bodyBytes, &requestData)
	if err != nil {
		return err
	}
	ucInput := usecaseinputports.TransferInput{
		BankAccountUserID: requestData.BankAccountUserID,
		FromBankAccountID: requestData.FromBankAccountID,
		ToBankAccountID:   requestData.ToBankAccountID,
		Amount:            requestData.Amount,
	}

	// Execute Use Case
	err = controller.usecase.Execute(ctx, ucInput)
	if err != nil {
		return err
	}
	return nil
}
