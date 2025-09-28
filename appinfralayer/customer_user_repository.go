package appinfralayer

import (
	"app/applogiclayer/domain/infrainterface"
	"app/applogiclayer/domain/model"
	"context"
)

type customerUserRepository struct {
	// e.g dynamodb client...
}

func NewCustomerUserRepository() infrainterface.CustomerUserRepository {
	return customerUserRepository{}
}

func (r customerUserRepository) Get(ctx context.Context, id model.CustomerUserID) (model.CustomerUser, error) {
	// fetch from db...
	dto := CustomerUserDTO{}
	return dto.toCustomerUserModel(), nil
}

type CustomerUserDTO struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	CreatedAt float64 `json:"createdAt"`
	UpdatedAt float64 `json:"updatedAt"`
}

func (dto CustomerUserDTO) toCustomerUserModel() model.CustomerUser {
	return model.CustomerUser{
		ID:        model.CustomerUserID(dto.ID),
		Name:      dto.Name,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}
}

func NewCustomerUserDTOByModel(m model.CustomerUser) CustomerUserDTO {
	return CustomerUserDTO{
		ID:        string(m.ID),
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
