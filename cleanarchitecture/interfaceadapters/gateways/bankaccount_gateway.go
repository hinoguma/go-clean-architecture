package gateways

import (
	"app/cleanarchitecture/enterprisebusinessrules/entities"
	"context"
)

type BankAccountGateway interface {
	Get(ctx context.Context, id string) (entities.BankAccount, error)
	Create(ctx context.Context, item entities.BankAccount) (string, error)
	Update(ctx context.Context, item entities.BankAccount) error
	Delete(ctx context.Context, id string) error
}
