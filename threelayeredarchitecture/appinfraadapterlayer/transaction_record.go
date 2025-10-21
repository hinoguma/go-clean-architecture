package appinfraadapterlayer

import (
	"app/threelayeredarchitecture/appinfralayer"
	"context"
)

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

type UpdateTransactionRecordRequestDTO struct {
	UpdateDatabaseItemRequestDTO
}

func NewUpdateTransactionRecordRequestDTO() *UpdateTransactionRecordRequestDTO {
	return &UpdateTransactionRecordRequestDTO{}
}

func (dto UpdateTransactionRecordRequestDTO) UpdateFieldRequests() appinfralayer.UpdateFieldRequests {
	var reqs appinfralayer.UpdateFieldRequests = make([]appinfralayer.UpdateFieldRequest, 0)
	if dto.UpdatedAt != nil {
		reqs.Append("updatedAt", dto.UpdatedAt)
	}
	return reqs
}

type TransactionRecordRepositoryAdapterIF interface {
	Get(ctx context.Context, id string) (TransactionRecordDTO, error)
	Create(ctx context.Context, itemDTO TransactionRecordDTO) error
	Update(ctx context.Context, id string, updateReq UpdateTransactionRecordRequestDTO) error
	Delete(ctx context.Context, id string) error
}

type TransactionRecordRepositoryAdapter struct {
	repository appinfralayer.TransactionRecordRepositoryIF
}

func NewTransactionRecordRepositoryAdapter(repository appinfralayer.TransactionRecordRepositoryIF) TransactionRecordRepositoryAdapterIF {
	return &TransactionRecordRepositoryAdapter{
		repository: repository,
	}
}

func (adapter *TransactionRecordRepositoryAdapter) Get(ctx context.Context, id string) (TransactionRecordDTO, error) {
	rawdata, err := adapter.repository.Get(ctx, id)
	if err != nil {
		return TransactionRecordDTO{}, err
	}
	dto := TransactionRecordDTO{}
	dto.SetFromRawData(rawdata)
	return dto, nil
}

func (adapter *TransactionRecordRepositoryAdapter) Create(ctx context.Context, itemDTO TransactionRecordDTO) error {
	return adapter.repository.Create(ctx, itemDTO.RawData())
}

func (adapter *TransactionRecordRepositoryAdapter) Update(ctx context.Context, id string, updateReq UpdateTransactionRecordRequestDTO) error {
	return adapter.repository.Update(ctx, id, updateReq.UpdateFieldRequests())
}

func (adapter *TransactionRecordRepositoryAdapter) Delete(ctx context.Context, id string) error {
	return adapter.repository.Delete(ctx, id)
}
