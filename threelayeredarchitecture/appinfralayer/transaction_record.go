package appinfralayer

type TransactionRecordRawData struct {
	ID string

	DatabaseItem
}

type TransactionRecordRepositoryIF interface {
	Get(id string) (TransactionRecordRawData, error)
	Create(item TransactionRecordRawData) error
	Update(item TransactionRecordRawData) error
	Delete(id string) error
}
