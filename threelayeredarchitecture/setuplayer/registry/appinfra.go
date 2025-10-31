package registry

import (
	"app/threelayeredarchitecture/appinfralayer"
	"database/sql"
)

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
