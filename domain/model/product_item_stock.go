package model

import (
	"app/crosscutting/utils"
)

type ProductItemStockID string

type ProductItemStock struct {
	ID                ProductItemStockID
	ProductItemInfoID ProductItemInfoID
	utils.DBItemCommonProps
}
