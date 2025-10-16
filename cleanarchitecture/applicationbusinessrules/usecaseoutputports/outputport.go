package usecaseoutputports

type OutputPort[UseCaseOutputType any, ViewModelType any] interface {
	ReceiveUseCaseOutput(output UseCaseOutputType) error
	Present() (ViewModelType, error)
}
