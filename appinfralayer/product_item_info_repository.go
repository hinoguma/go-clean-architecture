package appinfralayer

import (
	"app/applogiclayer/domain/infrainterface"
	"app/applogiclayer/domain/model"
	"context"
)

type productItemInfoRepository struct {
	// e.g dynamodb client...
}

func NewProductItemInfoRepository() infrainterface.ProductItemInfoRepository {
	return productItemInfoRepository{}
}

func (r productItemInfoRepository) Get(ctx context.Context, id model.ProductItemInfoID) (model.ProductItemInfo, error) {
	// fetch from db...
	dto := ProductItemInfoDTO{}
	return dto.toProductItemInfoModel(), nil
}

type ProductItemInfoDTO struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Currency  string  `json:"currency"`
	CreatedAt float64 `json:"createdAt"`
	UpdatedAt float64 `json:"updatedAt"`
}

func (dto ProductItemInfoDTO) toProductItemInfoModel() model.ProductItemInfo {
	return model.ProductItemInfo{
		ID:   model.ProductItemInfoID(dto.ID),
		Name: dto.Name,
		Price: model.Price{
			Amount:   dto.Price,
			Currency: model.Currency(dto.Currency),
		},
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}
}

func NewProductItemInfoDTOByModel(m model.ProductItemInfo) ProductItemInfoDTO {
	return ProductItemInfoDTO{
		ID:        string(m.ID),
		Name:      m.Name,
		Price:     m.Price.Amount,
		Currency:  string(m.Price.Currency),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
