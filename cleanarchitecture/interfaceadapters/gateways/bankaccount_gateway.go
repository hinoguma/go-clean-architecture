package gateways

import "app/cleanarchitecture/enterprisebusinessrules/entities"

type BankAccountGateway interface {
	Get(id string) (entities.BankAccount, error)
	Create(item entities.BankAccount) (string, error)
	Update(item entities.BankAccount) error
	Delete(id string) error
}
