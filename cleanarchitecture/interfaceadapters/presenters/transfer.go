package presenters

import (
	"app/cleanarchitecture/applicationbusinessrules/usecaseoutputports"
	"github.com/aws/aws-lambda-go/events"
	"net/http"
)

type TransferPresenter struct {
	store *TransferPresenterStore
}

func NewTransferPresenter(
	store *TransferPresenterStore,
) usecaseoutputports.TransferUseCaseOutputPort {
	return &TransferPresenter{
		store: store,
	}
}

func (presenter *TransferPresenter) ReceiveUseCaseOutput(output usecaseoutputports.TransferOutput) error {

	if output.ResultStatus == "Success" {
		presenter.store.Set(events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusOK,
			Body:       `{"transaction_id":"` + output.TransactionID + `"}`, // json body
		})
	}

	return nil
}

type TransferPresenterStore struct {
	response events.APIGatewayV2HTTPResponse
}

func (store TransferPresenterStore) Get() events.APIGatewayV2HTTPResponse {
	return store.response
}

func (store *TransferPresenterStore) Set(response events.APIGatewayV2HTTPResponse) {
	store.response = response
}

func NewTransferPresenterStore() *TransferPresenterStore {
	return &TransferPresenterStore{}
}
