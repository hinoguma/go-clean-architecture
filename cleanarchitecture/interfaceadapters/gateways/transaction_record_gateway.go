package gateways

import "app/cleanarchitecture/enterprisebusinessrules/entities"

type TransactionRecordGateway interface {
	Get(id string) (entities.TransactionRecord, error)
	Create(item entities.TransactionRecord) (string, error)
	Update(item entities.TransactionRecord) error
	Delete(id string) error
}
