package model

import (
	"app/crosscutting/utils"
)

type CustomerUserID string

type CustomerUser struct {
	ID   CustomerUserID `json:"id"`
	Name string         `json:"name"`
	utils.DBItemCommonProps
}
