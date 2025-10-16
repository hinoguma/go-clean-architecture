package appinfralayer

import (
	"app/experimentarchitecture/appinfrainterfacelayer"
	model2 "app/experimentarchitecture/applogiclayer/domain/model"
	"context"
)

type shoppingCartRepository struct {
	// e.g dynamodb client...
}

func NewShoppingCartRepository() appinfrainterfacelayer.ShoppingCartRepository {
	return shoppingCartRepository{}
}

func (r shoppingCartRepository) GetByUserID(ctx context.Context, userID model2.CustomerUserID) (model2.ShoppingCart, error) {
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

func (dto CartItemDTO) toCartItemModel() model2.CartItem {
	return model2.CartItem{
		ProductID: model2.ProductItemInfoID(dto.ProductID),
		Quantity:  dto.Quantity,
	}
}

func NewCartItemDTOByModel(m model2.CartItem) CartItemDTO {
	return CartItemDTO{
		ProductID: string(m.ProductID),
		Quantity:  m.Quantity,
	}
}

func (dto ShoppingCartDTO) toShoppingCartModel() model2.ShoppingCart {
	var items []model2.CartItem
	if len(dto.Items) > 0 {
		items = make([]model2.CartItem, len(dto.Items))
		for i, item := range dto.Items {
			items[i] = item.toCartItemModel()
		}
	}
	return model2.ShoppingCart{
		ID:        model2.ShoppingCartID(dto.ID),
		Items:     items,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}
}

func NewShoppingCartDTOByModel(m model2.ShoppingCart) ShoppingCartDTO {
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
