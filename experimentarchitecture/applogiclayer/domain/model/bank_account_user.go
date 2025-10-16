package model

import (
	"app/experimentarchitecture/crosscutting/utils"
	"context"
)

type BankUserID string

type BankUser struct {
	ID    BankUserID
	Name  string
	Email utils.Email
}

type BankUserAuthRequest struct {
	Ctx   context.Context
	Token string
}
