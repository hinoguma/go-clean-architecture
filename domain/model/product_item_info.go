package model

import (
	"app/crosscutting/utils"
)

type ProductItemInfoID string

type ProductItemInfo struct {
	ID    ProductItemInfoID `json:"id"`
	Name  string            `json:"name"`
	Price Price             `json:"price"`
	utils.DBItemCommonProps
}
