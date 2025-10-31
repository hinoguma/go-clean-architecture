package registry

import (
	"database/sql"
)

var RegistryAppInfra AppInfraRegistryIF
var RegistryAppInfraAdapter AppInfraAdapterIF
var RegistryDomainService DomainServiceRegistryIF
var RegistryUseCase UseCaseRegistryIF
var RegistryUserCallAdapter UserCallAdapterRegistryIF
var RegistryUserCallHandler UserCallHandlerRegistryIF

func Init(db *sql.DB) {
	InitAppInfraRegistry(db)
	InitAppInfraAdapterRegistry(RegistryAppInfra)
	InitDomainServiceRegistry(RegistryAppInfraAdapter)
	InitUseCaseRegistry(RegistryDomainService, RegistryAppInfraAdapter)
	InitUserCallAdapterRegistry(RegistryUseCase)
	InitUserCallHandlerRegistry(RegistryUserCallAdapter)
}

func InitAppInfraRegistry(db *sql.DB) {
	RegistryAppInfra = AppInfraRegistry{
		db: db,
	}
}

func InitAppInfraAdapterRegistry(appInfraRegistry AppInfraRegistryIF) {
	RegistryAppInfraAdapter = NewAppInfraAdapter(appInfraRegistry)
}

func InitDomainServiceRegistry(
	appInfraAdapterRegistry AppInfraAdapterIF,
) {
	RegistryDomainService = NewDomainServiceRegistry(appInfraAdapterRegistry)
}

func InitUseCaseRegistry(
	domainServiceRegistry DomainServiceRegistryIF,
	appInfraAdapterRegistry AppInfraAdapterIF,
) {
	RegistryUseCase = NewUseCaseRegistry(domainServiceRegistry, appInfraAdapterRegistry)
}

func InitUserCallAdapterRegistry(
	usecaseRegistry UseCaseRegistryIF,
) {
	RegistryUserCallAdapter = NewUserCallAdapterRegistry(usecaseRegistry)
}

func InitUserCallHandlerRegistry(
	userCallAdapterRegistry UserCallAdapterRegistryIF,
) {
	RegistryUserCallHandler = NewUserCallHandlerRegistry(userCallAdapterRegistry)
}
