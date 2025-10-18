package appinfraadapterlayer

import "app/threelayeredarchitecture/appinfralayer"

type TransactionRecordDTO struct {
	ID                string
	FromBankAccountID string
	ToBankAccountID   string
	Amount            int64
	Currency          string
	DatabaseItem
}

func (dto TransactionRecordDTO) RawData() appinfralayer.TransactionRecordRawData {
	rawdata := appinfralayer.TransactionRecordRawData{
		ID: dto.ID,
		DatabaseItem: appinfralayer.DatabaseItem{
			CreatedAt: dto.CreatedAt,
			UpdatedAt: dto.UpdatedAt,
		},
	}
	return rawdata
}

func (dto *TransactionRecordDTO) SetFromRawData(rawdata appinfralayer.TransactionRecordRawData) {
	dto.ID = rawdata.ID
	dto.CreatedAt = rawdata.CreatedAt
	dto.UpdatedAt = rawdata.UpdatedAt
}

type TransactionRecordRepositoryAdapterIF interface {
	Get(id string) (TransactionRecordDTO, error)
	Create(itemDTO TransactionRecordDTO) error
	Update(itemDTO TransactionRecordDTO) error
	Delete(id string) error
}

type TransactionRecordRepositoryAdapter struct {
	repository appinfralayer.TransactionRecordRepositoryIF
}

func NewTransactionRecordRepositoryAdapter(repository appinfralayer.TransactionRecordRepositoryIF) TransactionRecordRepositoryAdapterIF {
	return &TransactionRecordRepositoryAdapter{
		repository: repository,
	}
}

func (adapter *TransactionRecordRepositoryAdapter) Get(id string) (TransactionRecordDTO, error) {
	rawdata, err := adapter.repository.Get(id)
	if err != nil {
		return TransactionRecordDTO{}, err
	}
	dto := TransactionRecordDTO{}
	dto.SetFromRawData(rawdata)
	return dto, nil
}

func (adapter *TransactionRecordRepositoryAdapter) Create(itemDTO TransactionRecordDTO) error {
	return adapter.repository.Create(itemDTO.RawData())
}

func (adapter *TransactionRecordRepositoryAdapter) Update(itemDTO TransactionRecordDTO) error {
	return adapter.repository.Update(itemDTO.RawData())
}

func (adapter *TransactionRecordRepositoryAdapter) Delete(id string) error {
	return adapter.repository.Delete(id)
}
