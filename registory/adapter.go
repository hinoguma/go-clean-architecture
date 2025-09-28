package registory

import "app/adapter"

type AdapterFactory interface {
	NewAddItemToCartAdapter() adapter.AddItemToCartAdapter
}

type adapterFactory struct {
	useCaseFactory UseCaseFactory
}

func NewAdapterFactory(useCaseFactory UseCaseFactory) AdapterFactory {
	return adapterFactory{
		useCaseFactory: useCaseFactory,
	}
}

func (factory adapterFactory) NewAddItemToCartAdapter() adapter.AddItemToCartAdapter {
	return adapter.NewAddItemToCartAdapter(
		factory.useCaseFactory.NewAddItemToCartUseCase(),
	)
}
