package model

import (
	"app/threelayeredarchitecture/appinfraadapterlayer"
	"time"
)

type BankAccount struct {
	ID                  string
	OwnerBankCustomerID string
	Balance             Money

	DataItem
}

func (account *BankAccount) DTO() appinfraadapterlayer.BankAccountDTO {
	dto := appinfraadapterlayer.BankAccountDTO{
		ID:                  account.ID,
		OwnerBankCustomerID: account.OwnerBankCustomerID,
		BalanceAmount:       account.Balance.Amount,
		BalanceCurrency:     string(account.Balance.Currency),
	}
	dto.DatabaseItem.CreatedAt = account.CreatedAt.Unix()
	dto.DatabaseItem.UpdatedAt = account.UpdatedAt.Unix()
	return dto
}

func (model *BankAccount) Withdraw(money Money, now time.Time) {
	model.Balance.Amount -= money.Amount
	model.UpdatedAt = now
}

func (model *BankAccount) Deposit(money Money, now time.Time) {
	model.Balance.Amount += money.Amount
	model.UpdatedAt = now
}

func ConvertDTOToBankAccount(dto appinfraadapterlayer.BankAccountDTO) BankAccount {
	item := BankAccount{
		ID:                  dto.ID,
		OwnerBankCustomerID: dto.OwnerBankCustomerID,
		Balance: Money{
			Amount:   dto.BalanceAmount,
			Currency: Currency(dto.BalanceCurrency),
		},
	}
	item.DataItem.SetFromDTO(dto.DatabaseItem)
	return item
}
