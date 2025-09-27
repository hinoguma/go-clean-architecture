package model

type Price struct {
	Amount   float64  `json:"amount"`
	Currency Currency `json:"currency"`
}

type Currency string

const (
	JPY Currency = "JPY"
)
