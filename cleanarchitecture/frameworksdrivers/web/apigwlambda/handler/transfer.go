package handler

import (
	"app/cleanarchitecture/applicationbusinessrules/usecaseinteractors"
	"app/cleanarchitecture/frameworksdrivers/db"
	"app/cleanarchitecture/interfaceadapters/controllers"
	"app/cleanarchitecture/interfaceadapters/presenters"
	"context"
	"github.com/aws/aws-lambda-go/events"
)

func TransferHandler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	// set up presenter
	presentStore := presenters.NewTransferPresenterStore()
	presenter := presenters.NewTransferPresenter(presentStore)
	
	// set up gateways
	transactionRecordGW := db.NewTransactionRecordGateway()
	bankAccountGW := db.NewBankAccountGateway()

	// set up interactor
	interactor := usecaseinteractors.NewTransferUseCaseInteractor(
		presenter, transactionRecordGW, bankAccountGW,
	)

	// set up controller
	controller := controllers.NewTransferController(interactor)

	// call controller
	err := controller.Execute(ctx, event)
	if err != nil {
		return errorResponse(err)
	}
	// get response from presenter store
	resp := presentStore.Get()

	// return response
	return resp, nil
}
