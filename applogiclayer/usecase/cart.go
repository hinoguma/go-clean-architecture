package usecase

import "app/applogiclayer/domain/infrainterface"

type AddItemToCartRequest struct {
	CustomerUserID string
	ProductID      string
	Quantity       int
}

type AddItemToCartResult struct {
	CartID string
}

type AddItemToCartUseCase interface {
	Execute(input AddItemToCartRequest) AddItemToCartResult
}

func NewAddItemToCartUseCase(
	productItemInfoRepository infrainterface.ProductItemInfoRepository,
	customerUserRepository infrainterface.CustomerUserRepository,
	shoppingCartRepository infrainterface.ShoppingCartRepository,
) AddItemToCartUseCase {
	return addItemToCartUseCase{
		productItemInfoRepository: productItemInfoRepository,
		customerUserRepository:    customerUserRepository,
		shoppingCartRepository:    shoppingCartRepository,
	}
}

type addItemToCartUseCase struct {
	// auth service
	productItemInfoRepository infrainterface.ProductItemInfoRepository
	customerUserRepository    infrainterface.CustomerUserRepository
	shoppingCartRepository    infrainterface.ShoppingCartRepository
}

func (uc addItemToCartUseCase) Execute(req AddItemToCartRequest) AddItemToCartResult {
	return AddItemToCartResult{}
}
