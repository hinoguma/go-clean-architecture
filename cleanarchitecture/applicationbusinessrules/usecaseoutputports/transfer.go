package usecaseoutputports

type TransferOutput struct {
	ResultStatus  string
	TransactionID string
}

type TransferUseCaseOutputPort interface {
	ReceiveUseCaseOutput(output TransferOutput) error
}
