package usecaseinteractors

import (
	"app/cleanarchitecture/applicationbusinessrules/usecaseinputports"
	"app/cleanarchitecture/applicationbusinessrules/usecaseoutputports"
	"app/cleanarchitecture/interfaceadapters/gateways"
	"context"
)

type TransferUseCaseInteractor struct {
	outputport               usecaseoutputports.TransferUseCaseOutputPort
	transactionRecordGateway gateways.TransactionRecordGateway
	bankAccountGateway       gateways.BankAccountGateway
}

func NewTransferUseCaseInteractor(
	outputport usecaseoutputports.TransferUseCaseOutputPort,
	transactionRecordGateway gateways.TransactionRecordGateway,
	bankAccountGateway gateways.BankAccountGateway,
) usecaseinputports.TransferUseCaseInputPort {
	return &TransferUseCaseInteractor{
		outputport:               outputport,
		transactionRecordGateway: transactionRecordGateway,
		bankAccountGateway:       bankAccountGateway,
	}
}

func (interactor *TransferUseCaseInteractor) Execute(ctx context.Context, input usecaseinputports.TransferInput) error {

	// get bank accounts

	// create transaction record

	// check if balance is sufficient

	// transfer money

	// finish transfer

	// call output port
	output := usecaseoutputports.TransferOutput{}
	err := interactor.outputport.ReceiveUseCaseOutput(output)
	if err != nil {
		return err
	}
	return nil
}
