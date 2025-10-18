package appinfraadapterlayer

type TransactionRecordDTO struct {
	ID                string
	FromBankAccountID string
	ToBankAccountID   string
	Amount            int
	Currency          string
	DatabaseItem
}

type TransactionRecordRepositoryAdapterIF interface {
	Get(id string) (TransactionRecordDTO, error)
	Create(itemDTO TransactionRecordDTO) error
	Update(itemDTO TransactionRecordDTO) error
	Delete(id string) error
}
