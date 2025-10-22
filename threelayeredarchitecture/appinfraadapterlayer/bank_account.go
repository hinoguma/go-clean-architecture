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

func ConvertRawDataToBankAccountDTO(rawdata appinfralayer.BankAccountRawData) BankAccountDTO {
	item := BankAccountDTO{}
	item.ID = rawdata.ID
	item.CreatedAt = rawdata.CreatedAt
	item.UpdatedAt = rawdata.UpdatedAt
	return item
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
	Lock(ctx context.Context, id string, tx Transaction) (BankAccountDTO, error)
	TxCreate(ctx context.Context, itemDTO BankAccountDTO, tx Transaction) error
	TxUpdate(ctx context.Context, id string, updateReq UpdateBankAccountRequestDTO, tx Transaction) error
	TxDelete(ctx context.Context, id string, tx Transaction) error
}

func NewBankAccountRepositoryAdapter(
	repository appinfralayer.BankAccountRepositoryIF,
) BankAccountRepositoryAdapterIF {
	return &DatabaseItemRepositoryAdapter[
		appinfralayer.BankAccountRawData,
		BankAccountDTO,
		string,
		UpdateBankAccountRequestDTO,
	]{
		repository:              repository,
		convertRawDataToDtoFunc: ConvertRawDataToBankAccountDTO,
	}
}
