package registry

import (
	"app/threelayeredarchitecture/appinfraadapterlayer"
	"app/threelayeredarchitecture/appinfralayer"
	"app/threelayeredarchitecture/applogiclayer/usecase"
	"app/threelayeredarchitecture/applogiclayer/usecase/domain/service"
	"app/threelayeredarchitecture/usercalladapterlayer"
	"app/threelayeredarchitecture/usercallhandlerlayer"
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

/***********************************
 * app infra registry
 ***********************************/
type AppInfraRegistryIF interface {
	GetSQLDB() *sql.DB
	GetSQLClient() appinfralayer.SQLClient
	GetTxConnectionPool() appinfralayer.TxConnectionPoolIF
}

type AppInfraRegistry struct {
	db *sql.DB
}

func (registry AppInfraRegistry) GetSQLDB() *sql.DB {
	return registry.db
}

func (registry AppInfraRegistry) GetSQLClient() appinfralayer.SQLClient {
	return appinfralayer.NewPostgreSQLClient(registry.GetSQLDB())
}

func (registry AppInfraRegistry) GetTxConnectionPool() appinfralayer.TxConnectionPoolIF {
	return appinfralayer.NewTxConnectionManager()
}

/***********************************
 * app infra adapter registry
 ***********************************/
type AppInfraAdapterIF interface {
	GetTransactionRecordRepositoryAdapter() appinfraadapterlayer.TransactionRecordRepositoryAdapterIF
	GetBankAccountRepositoryAdapter() appinfraadapterlayer.BankAccountRepositoryAdapterIF
	GetBankCustomerRepositoryAdapter() appinfraadapterlayer.BankCustomerRepositoryAdapterIF
	GetTransactionManagerAdapter() appinfraadapterlayer.TransactionManagerAdapterIF
}

func NewAppInfraAdapter(
	infraRegistry AppInfraRegistryIF,
) AppInfraAdapterIF {
	return AppInfraAdapter{
		infraRegistry: infraRegistry,
	}
}

type AppInfraAdapter struct {
	infraRegistry AppInfraRegistryIF
}

func (adapter AppInfraAdapter) GetTransactionRecordRepositoryAdapter() appinfraadapterlayer.TransactionRecordRepositoryAdapterIF {
	return appinfraadapterlayer.NewTransactionRecordRepositoryAdapter(
		appinfralayer.NewTransactionRecordRepository(
			adapter.infraRegistry.GetSQLClient(),
			adapter.infraRegistry.GetTxConnectionPool(),
		),
	)
}

func (adapter AppInfraAdapter) GetBankAccountRepositoryAdapter() appinfraadapterlayer.BankAccountRepositoryAdapterIF {
	return appinfraadapterlayer.NewBankAccountRepositoryAdapter(
		appinfralayer.NewBankAccountRepository(
			adapter.infraRegistry.GetSQLClient(),
			adapter.infraRegistry.GetTxConnectionPool(),
		),
	)
}

func (adapter AppInfraAdapter) GetBankCustomerRepositoryAdapter() appinfraadapterlayer.BankCustomerRepositoryAdapterIF {
	return appinfraadapterlayer.NewBankCustomerRepositoryAdapter(
		appinfralayer.NewBankCustomerRepository(
			adapter.infraRegistry.GetSQLClient(),
			adapter.infraRegistry.GetTxConnectionPool(),
		),
	)
}

func (adapter AppInfraAdapter) GetTransactionManagerAdapter() appinfraadapterlayer.TransactionManagerAdapterIF {
	return appinfraadapterlayer.NewTransactionManagerAdapter(
		appinfralayer.NewTransactionManager(
			adapter.infraRegistry.GetSQLDB(),
		),
	)
}

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

/***********************************
 * usercall adapter registry
 ***********************************/
type UserCallAdapterRegistryIF interface {
	GetTransferAdapter() usercalladapterlayer.TransferAdapterIF
}

func NewUserCallAdapterRegistry(
	usecaseRegistry UseCaseRegistryIF,
) UserCallAdapterRegistryIF {
	return UserCallAdapterRegistry{
		usecaseRegistry: usecaseRegistry,
	}
}

type UserCallAdapterRegistry struct {
	usecaseRegistry UseCaseRegistryIF
}

func (registry UserCallAdapterRegistry) GetTransferAdapter() usercalladapterlayer.TransferAdapterIF {
	return usercalladapterlayer.NewTransferAdapter(
		registry.usecaseRegistry.GetTransferUseCase(),
	)
}

/***********************************
 * usercall handler registry
 ***********************************/
type UserCallHandlerRegistryIF interface {
	GetTransferHandler() usercallhandlerlayer.TransferHandlerIF
}

func NewUserCallHandlerRegistry(
	userCallAdapterRegistry UserCallAdapterRegistryIF,
) UserCallHandlerRegistryIF {
	return UserCallHandlerRegistry{
		userCallAdapterRegistry: userCallAdapterRegistry,
	}
}

type UserCallHandlerRegistry struct {
	userCallAdapterRegistry UserCallAdapterRegistryIF
}

func (registry UserCallHandlerRegistry) GetTransferHandler() usercallhandlerlayer.TransferHandlerIF {
	return usercallhandlerlayer.NewTransferHandler(
		registry.userCallAdapterRegistry.GetTransferAdapter(),
	)
}
