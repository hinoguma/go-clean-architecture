package appinfraadapterlayer

import "app/threelayeredarchitecture/appinfralayer"

type BankCustomerDTO struct {
	ID   string
	Name string
	DatabaseItem
}

func (dto BankCustomerDTO) RawData() appinfralayer.BankCustomerRawData {
	rawdata := appinfralayer.BankCustomerRawData{
		ID: dto.ID,
		DatabaseItem: appinfralayer.DatabaseItem{
			CreatedAt: dto.CreatedAt,
			UpdatedAt: dto.UpdatedAt,
		},
	}
	return rawdata
}

func (dto *BankCustomerDTO) SetFromRawData(rawdata appinfralayer.BankCustomerRawData) {
	dto.ID = rawdata.ID
	dto.CreatedAt = rawdata.CreatedAt
	dto.UpdatedAt = rawdata.UpdatedAt
}

type BankCustomerRepositoryAdapterIF interface {
	Get(id string) (BankCustomerDTO, error)
	Create(itemDTO BankCustomerDTO) error
	Update(itemDTO BankCustomerDTO) error
	Delete(id string) error
}

type BankCustomerRepositoryAdapter struct {
	repository appinfralayer.BankCustomerRepositoryIF
}

func NewBankCustomerRepositoryAdapter(repository appinfralayer.BankCustomerRepositoryIF) BankCustomerRepositoryAdapterIF {
	return &BankCustomerRepositoryAdapter{
		repository: repository,
	}
}

func (adapter *BankCustomerRepositoryAdapter) Get(id string) (BankCustomerDTO, error) {
	rawdata, err := adapter.repository.Get(id)
	if err != nil {
		return BankCustomerDTO{}, err
	}
	dto := BankCustomerDTO{}
	dto.SetFromRawData(rawdata)
	return dto, nil
}

func (adapter *BankCustomerRepositoryAdapter) Create(itemDTO BankCustomerDTO) error {
	return adapter.repository.Create(itemDTO.RawData())
}

func (adapter *BankCustomerRepositoryAdapter) Update(itemDTO BankCustomerDTO) error {
	return adapter.repository.Update(itemDTO.RawData())
}

func (adapter *BankCustomerRepositoryAdapter) Delete(id string) error {
	return adapter.repository.Delete(id)
}
