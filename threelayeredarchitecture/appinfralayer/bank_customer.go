package appinfralayer


type BankCustomerRawData struct {
	ID      string
	DatabaseItem
}



type BankCustomerRepositoryIF interface {
	Get(id string) (BankCustomerRawData, error)
	Create(item BankCustomerRawData) error
	Update(item BankCustomerRawData) error
	Delete(id string) error
}



