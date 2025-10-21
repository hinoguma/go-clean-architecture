package appinfraadapterlayer

import (
	"app/threelayeredarchitecture/appinfralayer"
	"context"
)

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

type UpdateBankCustomerRequestDTO struct {
	Name *string
	UpdateDatabaseItemRequestDTO
}

func NewUpdateBankCustomerRequestDTO() *UpdateBankCustomerRequestDTO {
	return &UpdateBankCustomerRequestDTO{}
}

func (dto *UpdateBankCustomerRequestDTO) WithName(val string) *UpdateBankCustomerRequestDTO {
	dto.Name = &val
	return dto
}

func (dto UpdateBankCustomerRequestDTO) UpdateFieldRequests() appinfralayer.UpdateFieldRequests {
	var reqs appinfralayer.UpdateFieldRequests = make([]appinfralayer.UpdateFieldRequest, 0)
	if dto.Name != nil {
		reqs.Append("name", dto.Name)
	}
	if dto.UpdatedAt != nil {
		reqs.Append("updatedAt", dto.UpdatedAt)
	}
	return reqs
}

type BankCustomerRepositoryAdapterIF interface {
	Get(ctx context.Context, id string) (BankCustomerDTO, error)
	Create(ctx context.Context, itemDTO BankCustomerDTO) error
	Update(ctx context.Context, id string, updateReq UpdateBankCustomerRequestDTO) error
	Delete(ctx context.Context, id string) error
}

type BankCustomerRepositoryAdapter struct {
	repository appinfralayer.BankCustomerRepositoryIF
}

func NewBankCustomerRepositoryAdapter(repository appinfralayer.BankCustomerRepositoryIF) BankCustomerRepositoryAdapterIF {
	return &BankCustomerRepositoryAdapter{
		repository: repository,
	}
}

func (adapter *BankCustomerRepositoryAdapter) Get(ctx context.Context, id string) (BankCustomerDTO, error) {
	rawdata, err := adapter.repository.Get(ctx, id)
	if err != nil {
		return BankCustomerDTO{}, err
	}
	dto := BankCustomerDTO{}
	dto.SetFromRawData(rawdata)
	return dto, nil
}

func (adapter *BankCustomerRepositoryAdapter) Create(ctx context.Context, itemDTO BankCustomerDTO) error {
	return adapter.repository.Create(ctx, itemDTO.RawData())
}

func (adapter *BankCustomerRepositoryAdapter) Update(ctx context.Context, id string, updateReq UpdateBankCustomerRequestDTO) error {
	return adapter.repository.Update(ctx, id, updateReq.UpdateFieldRequests())
}

func (adapter *BankCustomerRepositoryAdapter) Delete(ctx context.Context, id string) error {
	return adapter.repository.Delete(ctx, id)
}
