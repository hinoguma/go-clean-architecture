package db

import (
	"app/cleanarchitecture/enterprisebusinessrules/entities"
	"app/cleanarchitecture/interfaceadapters/gateways"
)

type transactionRecordGateway struct {
}

func NewTransactionRecordGateway() gateways.TransactionRecordGateway {
	return &transactionRecordGateway{}
}

func (t transactionRecordGateway) Get(id string) (entities.TransactionRecord, error) {
	//TODO implement me
	panic("implement me")
}

func (t transactionRecordGateway) Create(item entities.TransactionRecord) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (t transactionRecordGateway) Update(item entities.TransactionRecord) error {
	//TODO implement me
	panic("implement me")
}

func (t transactionRecordGateway) Delete(id string) error {
	//TODO implement me
	panic("implement me")
}
