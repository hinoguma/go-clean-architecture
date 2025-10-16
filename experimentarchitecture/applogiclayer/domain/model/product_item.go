package model

type Currency string

type Price struct {
	Amount   float64  `json:"amount"`
	Currency Currency `json:"currency"`
}

type ProductItemInfoID string

type ProductItemInfo struct {
	ID        ProductItemInfoID `json:"id"`
	Name      string            `json:"name"`
	Price     Price             `json:"price"`
	CreatedAt float64           `json:"createdAt"`
	UpdatedAt float64           `json:"updatedAt"`
}
