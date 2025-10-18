package model

type Currency string

const (
	USD Currency = "USD"
	EUR Currency = "EUR"
	JPY Currency = "JPY"
)

type Money struct {
	Amount   int64
	Currency Currency
}

func (value Money) IsValid() bool {
	switch value.Currency {
	case USD, EUR, JPY:
	default:
		return false
	}
	return value.Amount >= 0
}


func NewUSD(amount int64) Money {
	return Money{
		Amount:   amount,
		Currency: USD,
	}
}

func NewEUR(amount int64) Money {
	return Money{
		Amount:   amount,
		Currency: EUR,
	}
}


func NewJPY(amount int64) Money {
	return Money{
		Amount:   amount,
		Currency: JPY,
	}
}