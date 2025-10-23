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

func convertRawDataToBankCustomerDTO(rawdata appinfralayer.BankCustomerRawData) BankCustomerDTO {
	item := BankCustomerDTO{}
	item.ID = rawdata.ID
	item.CreatedAt = rawdata.CreatedAt
	item.UpdatedAt = rawdata.UpdatedAt
	return item
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

type BankAccountAuthResultDTO struct {
	Success  bool
	Customer BankCustomerDTO

	Err         error
	ErrorReason AuthErrorReason
}

type AuthErrorReason struct {
	TokenExpired  bool
	InvalidToken  bool
	InternalError bool
}

type BankCustomerRepositoryAdapterIF interface {
	Authenticate(ctx context.Context, req AuthenticateRequestDTO) BankAccountAuthResultDTO
	Get(ctx context.Context, id string) (BankCustomerDTO, error)
	Create(ctx context.Context, itemDTO BankCustomerDTO) error
	Update(ctx context.Context, id string, updateReq UpdateBankCustomerRequestDTO) error
	Delete(ctx context.Context, id string) error
}

type BankCustomerRepositoryAdapter struct {
	// db
	DatabaseItemRepositoryAdapter[
		appinfralayer.BankCustomerRawData,
		BankCustomerDTO,
		string,
		UpdateBankCustomerRequestDTO,
	]

	// auth
}

func NewBankCustomerRepositoryAdapter(repository appinfralayer.BankCustomerRepositoryIF) BankCustomerRepositoryAdapterIF {
	adapter := BankCustomerRepositoryAdapter{}
	adapter.repository = repository
	adapter.convertRawDataToDtoFunc = convertRawDataToBankCustomerDTO
	return &adapter
}

func (adapter BankCustomerRepositoryAdapter) Authenticate(ctx context.Context, req AuthenticateRequestDTO) BankAccountAuthResultDTO {
	// todo: implement auth client in app infla layer

	return BankAccountAuthResultDTO{
		Success:     false,
		Customer:    BankCustomerDTO{},
		Err:         nil,
		ErrorReason: AuthErrorReason{},
	}
}
