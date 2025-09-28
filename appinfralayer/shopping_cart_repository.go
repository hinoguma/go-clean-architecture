package appinfralayer

import (
	"app/applogiclayer/domain/infrainterface"
	"app/applogiclayer/domain/model"
	"context"
)

type shoppingCartRepository struct {
	// e.g dynamodb client...
}

func NewShoppingCartRepository() infrainterface.ShoppingCartRepository {
	return shoppingCartRepository{}
}

func (r shoppingCartRepository) GetByUserID(ctx context.Context, userID model.CustomerUserID) (model.ShoppingCart, error) {
	// fetch from db...
	dto := ShoppingCartDTO{}
	return dto.toShoppingCartModel(), nil
}

type ShoppingCartDTO struct {
	ID             string        `json:"id"`
	CustomerUserID string        `json:"customerUserId"`
	Items          []CartItemDTO `json:"items"`
	CreatedAt      float64       `json:"createdAt"`
	UpdatedAt      float64       `json:"updatedAt"`
}

type CartItemDTO struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

func (dto CartItemDTO) toCartItemModel() model.CartItem {
	return model.CartItem{
		ProductID: model.ProductItemInfoID(dto.ProductID),
		Quantity:  dto.Quantity,
	}
}

func NewCartItemDTOByModel(m model.CartItem) CartItemDTO {
	return CartItemDTO{
		ProductID: string(m.ProductID),
		Quantity:  m.Quantity,
	}
}

func (dto ShoppingCartDTO) toShoppingCartModel() model.ShoppingCart {
	var items []model.CartItem
	if len(dto.Items) > 0 {
		items = make([]model.CartItem, len(dto.Items))
		for i, item := range dto.Items {
			items[i] = item.toCartItemModel()
		}
	}
	return model.ShoppingCart{
		ID:        model.ShoppingCartID(dto.ID),
		Items:     items,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}
}

func NewShoppingCartDTOByModel(m model.ShoppingCart) ShoppingCartDTO {
	var items []CartItemDTO
	if len(m.Items) > 0 {
		items := make([]CartItemDTO, len(m.Items))
		for i, item := range m.Items {
			items[i] = NewCartItemDTOByModel(item)
		}
	}
	return ShoppingCartDTO{
		ID:        string(m.ID),
		Items:     items,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
