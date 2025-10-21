package appinfraadapterlayer

import (
	"app/threelayeredarchitecture/appinfralayer"
	"context"
)

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

type UpdateBankAccountRequestDTO struct {
	OwnerBankCustomerID *string
	BalanceAmount       *int64
	BalanceCurrency     *string
	UpdateDatabaseItemRequestDTO
}

func NewUpdateBankAccountRequestDTO() *UpdateBankAccountRequestDTO {
	return &UpdateBankAccountRequestDTO{}
}

func (dto *UpdateBankAccountRequestDTO) WithOwnerBankCustomerID(val string) *UpdateBankAccountRequestDTO {
	dto.OwnerBankCustomerID = &val
	return dto
}

func (dto *UpdateBankAccountRequestDTO) WithBalanceAmount(val int64) *UpdateBankAccountRequestDTO {
	dto.BalanceAmount = &val
	return dto
}

func (dto *UpdateBankAccountRequestDTO) WithBalanceCurrency(val string) *UpdateBankAccountRequestDTO {
	dto.BalanceCurrency = &val
	return dto
}

func (dto UpdateBankAccountRequestDTO) UpdateFieldRequests() appinfralayer.UpdateFieldRequests {
	var reqs appinfralayer.UpdateFieldRequests = make([]appinfralayer.UpdateFieldRequest, 0)
	if dto.OwnerBankCustomerID != nil {
		reqs.Append("ownerBankCustomerId", dto.OwnerBankCustomerID)
	}
	if dto.BalanceAmount != nil {
		reqs.Append("balanceAmount", dto.BalanceAmount)
	}
	if dto.BalanceCurrency != nil {
		reqs.Append("balanceCurrency", dto.BalanceCurrency)
	}
	return reqs
}

type BankAccountRepositoryAdapterIF interface {
	Get(ctx context.Context, id string) (BankAccountDTO, error)
	Create(ctx context.Context, itemDTO BankAccountDTO) error
	Update(ctx context.Context, id string, updateReq UpdateBankAccountRequestDTO) error
	Delete(ctx context.Context, id string) error
}

type BankAccountRepositoryAdapter struct {
	repository appinfralayer.BankAccountRepositoryIF
}

func NewBankAccountRepositoryAdapter(repository appinfralayer.BankAccountRepositoryIF) BankAccountRepositoryAdapterIF {
	return &BankAccountRepositoryAdapter{
		repository: repository,
	}
}

func (adapter *BankAccountRepositoryAdapter) Get(ctx context.Context, id string) (BankAccountDTO, error) {
	rawdata, err := adapter.repository.Get(ctx, id)
	if err != nil {
		return BankAccountDTO{}, err
	}
	dto := BankAccountDTO{}
	dto.SetFromRawData(rawdata)
	return dto, nil
}

func (adapter *BankAccountRepositoryAdapter) Create(ctx context.Context, itemDTO BankAccountDTO) error {
	return adapter.repository.Create(ctx, itemDTO.RawData())
}

func (adapter *BankAccountRepositoryAdapter) Update(ctx context.Context, id string, updateReq UpdateBankAccountRequestDTO) error {
	return adapter.repository.Update(ctx, id, updateReq.UpdateFieldRequests())
}

func (adapter *BankAccountRepositoryAdapter) Delete(ctx context.Context, id string) error {
	return adapter.repository.Delete(ctx, id)
}
