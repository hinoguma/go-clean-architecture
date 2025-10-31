package appinfraadapterlayer

import (
	"app/threelayeredarchitecture/appinfralayer"
	"app/threelayeredarchitecture/crosscuttinglayer"
	"context"
)

type DatabaseItem struct {
	CreatedAt int64
	UpdatedAt int64
}

type DatabaseItemIF[
	RawDataType appinfralayer.SQLDatabaseItem[RawDataType],
] interface {
	RawData() RawDataType
}

type UpdateDatabaseItemRequestIF interface {
	UpdateFieldRequests() appinfralayer.UpdateFieldRequests
}

type UpdateDatabaseItemRequestDTO struct {
	UpdatedAt *int64
}

func (dto *UpdateDatabaseItemRequestDTO) WithUpdatedAt(val int64) *UpdateDatabaseItemRequestDTO {
	dto.UpdatedAt = &val
	return dto
}

type DatabaseItemRepositoryAdapterIF[
	RawDataType appinfralayer.SQLDatabaseItem[RawDataType],
	DTOType DatabaseItemIF[RawDataType],
	idType string | int64,
	updateRequestType UpdateDatabaseItemRequestIF,
] interface {
	Get(ctx context.Context, id idType) (DTOType, error)
	Create(ctx context.Context, itemDTO DTOType) error
	Put(ctx context.Context, itemDTO DTOType) error
	Update(ctx context.Context, id idType, updateReq updateRequestType) error
	Delete(ctx context.Context, id idType) error
	Lock(ctx context.Context, id idType, tx Transaction) (DTOType, error)
	TxCreate(ctx context.Context, itemDTO DTOType, tx Transaction) error
	TxPut(ctx context.Context, itemDTO DTOType, tx Transaction) error
	TxUpdate(ctx context.Context, id idType, updateReq updateRequestType, tx Transaction) error
	TxDelete(ctx context.Context, id idType, tx Transaction) error
}

type DatabaseItemRepositoryAdapter[
	RawDataType appinfralayer.SQLDatabaseItem[RawDataType],
	DTOType DatabaseItemIF[RawDataType],
	idType string | int64,
	updateRequestType UpdateDatabaseItemRequestIF,
] struct {
	repository              appinfralayer.TableRepositoryIF[RawDataType, idType]
	convertRawDataToDtoFunc func(raw RawDataType) DTOType
}

func (adapter *DatabaseItemRepositoryAdapter[RawDataType, DTOType, idType, updateReqType]) Get(ctx context.Context, id idType) (DTOType, error) {
	var dto DTOType
	rawdata, err := adapter.repository.Get(ctx, id)
	if err != nil {
		return dto, crosscuttinglayer.ErrLift(err, ctx)
	}
	dto = adapter.convertRawDataToDtoFunc(rawdata)
	return dto, nil
}

func (adapter *DatabaseItemRepositoryAdapter[RawDataType, DTOType, idType, updateReqType]) Create(ctx context.Context, itemDTO DTOType) error {
	return adapter.repository.Create(ctx, itemDTO.RawData())
}

func (adapter *DatabaseItemRepositoryAdapter[RawDataType, DTOType, idType, updateReqType]) Put(ctx context.Context, itemDTO DTOType) error {
	return adapter.repository.Put(ctx, itemDTO.RawData())
}

func (adapter *DatabaseItemRepositoryAdapter[RawDataType, DTOType, idType, updateReqType]) Update(ctx context.Context, id idType, updateReq updateReqType) error {
	return adapter.repository.Update(ctx, id, updateReq.UpdateFieldRequests())
}

func (adapter *DatabaseItemRepositoryAdapter[RawDataType, DTOType, idType, updateReqType]) Delete(ctx context.Context, id idType) error {
	return adapter.repository.Delete(ctx, id)
}

func (adapter *DatabaseItemRepositoryAdapter[RawDataType, DTOType, idType, updateReqType]) Lock(ctx context.Context, id idType, tx Transaction) (DTOType, error) {
	var dto DTOType
	rawdata, err := adapter.repository.Lock(ctx, id, tx.ID)
	if err != nil {
		return dto, crosscuttinglayer.ErrLift(err, ctx)
	}
	dto = adapter.convertRawDataToDtoFunc(rawdata)
	return dto, nil
}

func (adapter *DatabaseItemRepositoryAdapter[RawDataType, DTOType, idType, updateReqType]) TxCreate(ctx context.Context, itemDTO DTOType, tx Transaction) error {
	return adapter.repository.TxCreate(ctx, itemDTO.RawData(), tx.ID)
}

func (adapter *DatabaseItemRepositoryAdapter[RawDataType, DTOType, idType, updateReqType]) TxPut(ctx context.Context, itemDTO DTOType, tx Transaction) error {
	return adapter.repository.TxPut(ctx, itemDTO.RawData(), tx.ID)
}

func (adapter *DatabaseItemRepositoryAdapter[RawDataType, DTOType, idType, updateReqType]) TxUpdate(ctx context.Context, id idType, updateReq updateReqType, tx Transaction) error {
	return adapter.repository.TxUpdate(ctx, id, updateReq.UpdateFieldRequests(), tx.ID)
}

func (adapter *DatabaseItemRepositoryAdapter[RawDataType, DTOType, idType, updateReqType]) TxDelete(ctx context.Context, id idType, tx Transaction) error {
	return adapter.repository.TxDelete(ctx, id, tx.ID)
}
