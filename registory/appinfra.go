package registory

import (
	"app/appinfralayer"
	"app/applogiclayer/domain/infrainterface"
)

type InfraFactory interface {
	NewProductItemInfoRepository() infrainterface.ProductItemInfoRepository
	NewShoppingCartRepository() infrainterface.ShoppingCartRepository
	NewCustomerUserRepository() infrainterface.CustomerUserRepository
}

type infraFactory struct {
	// DBClient        utils.DBClient
}

func NewInfraFactory() InfraFactory {
	return infraFactory{}
}

func (factory infraFactory) NewProductItemInfoRepository() infrainterface.ProductItemInfoRepository {
	return appinfralayer.NewProductItemInfoRepository()
}

func (factory infraFactory) NewShoppingCartRepository() infrainterface.ShoppingCartRepository {
	return appinfralayer.NewShoppingCartRepository()
}
func (factory infraFactory) NewCustomerUserRepository() infrainterface.CustomerUserRepository {
	return appinfralayer.NewCustomerUserRepository()
}
