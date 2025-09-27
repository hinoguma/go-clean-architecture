package model

import (
	"app/crosscutting/utils"
)

type ShoppingCartID string

type ShoppingCartItem struct {
	ProductItemInfoID ProductItemInfoID
	Quantity          int
	AddRecordID       ShoppingCartOperationRecordID
}

type ShoppingCart struct {
	ID             ShoppingCartID     `json:"id"`
	CustomerUserID CustomerUserID     `json:"customerUserId"`
	Items          []ShoppingCartItem `json:"items"`
	DataItemCommonProps
}

type ShoppingCartOperationRecordID string

type CartOperationType string

const (
	CartOperationAdd     CartOperationType = "add"
	CartOperationCountUp CartOperationType = "count_up"
	CartOperationRemove  CartOperationType = "remove"
)

type ShoppingCartOperationRecord struct {
	ID                ShoppingCartOperationRecordID `json:"id"`
	ShoppingCartID    ShoppingCartID                `json:"shoppingCartId"`
	OperationType     CartOperationType             `json:"operationType"`
	ProductItemInfoID ProductItemInfoID             `json:"productItemInfoId"`
	utils.DBItemCommonProps
}
