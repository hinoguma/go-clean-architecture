package db

import (
	"app/cleanarchitecture/enterprisebusinessrules/entities"
	"app/cleanarchitecture/interfaceadapters/gateways"
)

type bankAccountGateway struct {
}

func NewBankAccountGateway() gateways.BankAccountGateway {
	return &bankAccountGateway{}
}

func (t bankAccountGateway) Get(id string) (entities.BankAccount, error) {
	//TODO implement me
	panic("implement me")
}

func (t bankAccountGateway) Create(item entities.BankAccount) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (t bankAccountGateway) Update(item entities.BankAccount) error {
	//TODO implement me
	panic("implement me")
}

func (t bankAccountGateway) Delete(id string) error {
	//TODO implement me
	panic("implement me")
}
