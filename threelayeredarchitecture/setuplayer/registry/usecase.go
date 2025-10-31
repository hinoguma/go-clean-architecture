package registry

import (
	"app/threelayeredarchitecture/applogiclayer/usecase"
)

/***********************************
 * usecase registry
 ***********************************/
type UseCaseRegistryIF interface {
	GetTransferUseCase() usecase.TransferUsecaseIF
}

func NewUseCaseRegistry(
	domainServiceRegistry DomainServiceRegistryIF,
	appInfraAdapterRegistry AppInfraAdapterIF,
) UseCaseRegistryIF {
	return UseCaseRegistry{
		domainServiceRegistry:   domainServiceRegistry,
		appInfraAdapterRegistry: appInfraAdapterRegistry,
	}
}

type UseCaseRegistry struct {
	domainServiceRegistry   DomainServiceRegistryIF
	appInfraAdapterRegistry AppInfraAdapterIF
}

func (registry UseCaseRegistry) GetTransferUseCase() usecase.TransferUsecaseIF {
	return usecase.NewTransferUsecase(
		registry.domainServiceRegistry.GetBankCustomerAuthService(),
		registry.domainServiceRegistry.GetTransferService(),
	)
}
