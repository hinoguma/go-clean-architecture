package registry

import (
	"app/threelayeredarchitecture/appinfraadapterlayer"
	"app/threelayeredarchitecture/appinfralayer"
)

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
