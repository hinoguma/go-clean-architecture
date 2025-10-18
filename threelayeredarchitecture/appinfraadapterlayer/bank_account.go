package appinfraadapterlayer



type BankAccountDTO struct {
	ID      string
	OwnerBankCustomerID string
	Balance int
	DatabaseItem
}

type BankAccountRepositoryAdapterIF interface {
	Get(id string) (BankAccountDTO, error)
	Create(itemDTO BankAccountDTO) error
	Update(itemDTO BankAccountDTO) error
	Delete(id string) error
}
