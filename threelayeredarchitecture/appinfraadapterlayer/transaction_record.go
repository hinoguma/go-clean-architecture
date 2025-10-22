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

func convertRawDataToBankAccountDTO(rawdata appinfralayer.TransactionRecordRawData) TransactionRecordDTO {
	item := TransactionRecordDTO{}
	item.ID = rawdata.ID
	item.CreatedAt = rawdata.CreatedAt
	item.UpdatedAt = rawdata.UpdatedAt
	return item
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

func NewTransactionRecordRepositoryAdapter(repository appinfralayer.TransactionRecordRepositoryIF) TransactionRecordRepositoryAdapterIF {
	return &DatabaseItemRepositoryAdapter[
		appinfralayer.TransactionRecordRawData,
		TransactionRecordDTO,
		string,
		UpdateTransactionRecordRequestDTO,
	]{
		repository:              repository,
		convertRawDataToDtoFunc: convertRawDataToBankAccountDTO,
	}
}
