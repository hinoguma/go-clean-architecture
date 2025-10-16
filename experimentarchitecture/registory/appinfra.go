package registory

import (
	"app/experimentarchitecture/appinfrainterfacelayer"
	appinfralayer2 "app/experimentarchitecture/appinfralayer"
)

type InfraFactory interface {
	NewProductItemInfoRepository() appinfrainterfacelayer.ProductItemInfoRepository
	NewShoppingCartRepository() appinfrainterfacelayer.ShoppingCartRepository
	NewCustomerUserRepository() appinfrainterfacelayer.CustomerUserRepository
}

type infraFactory struct {
	// DBClient        utils.DBClient
}

func NewInfraFactory() InfraFactory {
	return infraFactory{}
}

func (factory infraFactory) NewProductItemInfoRepository() appinfrainterfacelayer.ProductItemInfoRepository {
	return appinfralayer2.NewProductItemInfoRepository()
}

func (factory infraFactory) NewShoppingCartRepository() appinfrainterfacelayer.ShoppingCartRepository {
	return appinfralayer2.NewShoppingCartRepository()
}
func (factory infraFactory) NewCustomerUserRepository() appinfrainterfacelayer.CustomerUserRepository {
	return appinfralayer2.NewCustomerUserRepository()
}
