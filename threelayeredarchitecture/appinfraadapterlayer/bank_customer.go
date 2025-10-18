package appinfraadapterlayer



type BankCustomerDTO struct {
	ID      string
	Name    string
	DatabaseItem
}

type BankCustomerRepositoryAdapterIF interface {
	Get(id string) (BankCustomerDTO, error)
	Create(itemDTO BankCustomerDTO) error
	Update(itemDTO BankCustomerDTO) error
	Delete(id string) error
}
