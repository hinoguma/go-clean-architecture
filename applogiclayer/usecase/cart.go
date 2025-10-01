package usecase

import (
	"app/applogiclayer/domain/infrainterface"
	"app/applogiclayer/domain/service"
)

type AddItemToCartRequest struct {
	CustomerUserID string
	ProductID      string
	Quantity       int
	BaseUseCaseRequest
}

type AddItemToCartResult struct {
	CartID string
	BaseUseCaseResult
}

type AddItemToCartUseCase interface {
	Execute(input AddItemToCartRequest) (AddItemToCartResult, error)
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
	bankUserAuthService       service.BankUserAuthService
	productItemInfoRepository infrainterface.ProductItemInfoRepository
	customerUserRepository    infrainterface.CustomerUserRepository
	shoppingCartRepository    infrainterface.ShoppingCartRepository
}

func (uc addItemToCartUseCase) Execute(req AddItemToCartRequest) (AddItemToCartResult, error) {

	// Authenticate
	user, err := uc.bankUserAuthService.Authenticate(
		req.BankAccountUserAuthRequest(),
	)
	if err != nil {
		return AddItemToCartResult{}, NewAuthorizationError(err)
	}

	// Validation

	// Business Logic

	// response

	return AddItemToCartResult{}, nil
}
