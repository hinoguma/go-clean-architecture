package appinfraadapterlayer

import (
	"app/threelayeredarchitecture/appinfralayer"
	"app/threelayeredarchitecture/crosscuttinglayer"
	"context"
)

type TransactionRecordDTO struct {
	ID                string
	Type              string
	IdempotencyKey    string
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

func convertRawDataToTransactionRecordDTO(rawdata appinfralayer.TransactionRecordRawData) TransactionRecordDTO {
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
	GetByIdempotencyKey(ctx context.Context, key string) (TransactionRecordDTO, error)
	Create(ctx context.Context, itemDTO TransactionRecordDTO) error
	Put(ctx context.Context, itemDTO TransactionRecordDTO) error
	Update(ctx context.Context, id string, updateReq UpdateTransactionRecordRequestDTO) error
	Delete(ctx context.Context, id string) error
	TxCreate(ctx context.Context, itemDTO TransactionRecordDTO, tx Transaction) error
	TxPut(ctx context.Context, itemDTO TransactionRecordDTO, tx Transaction) error
}

func NewTransactionRecordRepositoryAdapter(repository appinfralayer.TransactionRecordRepositoryIF) TransactionRecordRepositoryAdapterIF {
	adapter := TransactionRecordRepositoryAdapter{}
	adapter.repository = repository
	adapter.convertRawDataToDtoFunc = convertRawDataToTransactionRecordDTO
	return &adapter
}

type TransactionRecordRepositoryAdapter struct {
	DatabaseItemRepositoryAdapter[
		appinfralayer.TransactionRecordRawData,
		TransactionRecordDTO,
		string,
		UpdateTransactionRecordRequestDTO,
	]
	repository appinfralayer.TransactionRecordRepositoryIF
}

func (adapter TransactionRecordRepositoryAdapter) GetByIdempotencyKey(ctx context.Context, key string) (TransactionRecordDTO, error) {
	var dto TransactionRecordDTO
	rawdata, err := adapter.repository.GetByIdempotencyKey(ctx, key)
	if err != nil {
		return dto, crosscuttinglayer.ErrLift(err, ctx)
	}
	dto = adapter.convertRawDataToDtoFunc(rawdata)
	return dto, nil
}
