package appinfralayer


type BankAccountRawData struct {
	ID      string


	DatabaseItem
}



type BankAccountRepositoryIF interface {
	Get(id string) (BankAccountRawData, error)
	Create(item BankAccountRawData) error
	Update(item BankAccountRawData) error
	Delete(id string) error
}



