package registory

import (
	"app/experimentarchitecture/applogiclayer/usecase"
)

type UseCaseFactory interface {
	NewAddItemToCartUseCase() usecase.AddItemToCartUseCase
}

func NewUseCaseFactory(infraFactory InfraFactory) UseCaseFactory {
	return useCaseFactory{
		infraFactory: infraFactory,
	}
}

type useCaseFactory struct {
	infraFactory InfraFactory
}

func (factory useCaseFactory) NewAddItemToCartUseCase() usecase.AddItemToCartUseCase {
	return usecase.NewAddItemToCartUseCase(
		factory.infraFactory.NewProductItemInfoRepository(),
		factory.infraFactory.NewCustomerUserRepository(),
		factory.infraFactory.NewShoppingCartRepository(),
	)
}
