package usecase

import (
	"app/experimentarchitecture/appinfrainterfacelayer"
	"app/experimentarchitecture/applogiclayer/domain/service"
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
	productItemInfoRepository appinfrainterfacelayer.ProductItemInfoRepository,
	customerUserRepository appinfrainterfacelayer.CustomerUserRepository,
	shoppingCartRepository appinfrainterfacelayer.ShoppingCartRepository,
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
	productItemInfoRepository appinfrainterfacelayer.ProductItemInfoRepository
	customerUserRepository    appinfrainterfacelayer.CustomerUserRepository
	shoppingCartRepository    appinfrainterfacelayer.ShoppingCartRepository
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
