package appinfraadapterlayer

import "app/threelayeredarchitecture/appinfralayer"

type BankAccountDTO struct {
	ID                  string
	OwnerBankCustomerID string
	BalanceAmount       int64
	BalanceCurrency     string
	DatabaseItem
}

func (dto BankAccountDTO) RawData() appinfralayer.BankAccountRawData {
	rawdata := appinfralayer.BankAccountRawData{
		ID: dto.ID,
		DatabaseItem: appinfralayer.DatabaseItem{
			CreatedAt: dto.CreatedAt,
			UpdatedAt: dto.UpdatedAt,
		},
	}
	return rawdata
}

func (dto *BankAccountDTO) SetFromRawData(rawdata appinfralayer.BankAccountRawData) {
	dto.ID = rawdata.ID
	dto.CreatedAt = rawdata.CreatedAt
	dto.UpdatedAt = rawdata.UpdatedAt
}

type BankAccountRepositoryAdapterIF interface {
	Get(id string) (BankAccountDTO, error)
	Create(itemDTO BankAccountDTO) error
	Update(itemDTO BankAccountDTO) error
	Delete(id string) error
}

type BankAccountRepositoryAdapter struct {
	repository appinfralayer.BankAccountRepositoryIF
}

func NewBankAccountRepositoryAdapter(repository appinfralayer.BankAccountRepositoryIF) BankAccountRepositoryAdapterIF {
	return &BankAccountRepositoryAdapter{
		repository: repository,
	}
}

func (adapter *BankAccountRepositoryAdapter) Get(id string) (BankAccountDTO, error) {
	rawdata, err := adapter.repository.Get(id)
	if err != nil {
		return BankAccountDTO{}, err
	}
	dto := BankAccountDTO{}
	dto.SetFromRawData(rawdata)
	return dto, nil
}

func (adapter *BankAccountRepositoryAdapter) Create(itemDTO BankAccountDTO) error {
	return adapter.repository.Create(itemDTO.RawData())
}

func (adapter *BankAccountRepositoryAdapter) Update(itemDTO BankAccountDTO) error {
	return adapter.repository.Update(itemDTO.RawData())
}

func (adapter *BankAccountRepositoryAdapter) Delete(id string) error {
	return adapter.repository.Delete(id)
}
