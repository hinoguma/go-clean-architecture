package registry

import (
	"app/threelayeredarchitecture/applogiclayer/usecase/domain/service"
)

/***********************************
 * domain service registry
 ***********************************/
type DomainServiceRegistryIF interface {
	GetBankCustomerAuthService() service.BankCustomerAuthServiceIF
	GetTransferService() service.TransferServiceIF
	GetBankAccountLockService() service.BankAccountLockServiceIF
}

func NewDomainServiceRegistry(
	infraAdapterRegistry AppInfraAdapterIF,
) DomainServiceRegistryIF {
	return DomainServiceRegistry{
		infraAdapterRegistry: infraAdapterRegistry,
	}
}

type DomainServiceRegistry struct {
	infraAdapterRegistry AppInfraAdapterIF
}

func (registry DomainServiceRegistry) GetBankCustomerAuthService() service.BankCustomerAuthServiceIF {
	return service.NewBankCustomerAuthService(
		registry.infraAdapterRegistry.GetBankCustomerRepositoryAdapter(),
	)
}

func (registry DomainServiceRegistry) GetBankAccountLockService() service.BankAccountLockServiceIF {
	return service.NewBankAccountLockService(
		registry.infraAdapterRegistry.GetBankAccountRepositoryAdapter(),
	)
}

func (registry DomainServiceRegistry) GetTransferService() service.TransferServiceIF {
	return service.NewTransferService(
		registry.GetBankAccountLockService(),
		registry.infraAdapterRegistry.GetTransactionManagerAdapter(),
		registry.infraAdapterRegistry.GetBankCustomerRepositoryAdapter(),
		registry.infraAdapterRegistry.GetBankAccountRepositoryAdapter(),
		registry.infraAdapterRegistry.GetTransactionRecordRepositoryAdapter(),
	)
}
