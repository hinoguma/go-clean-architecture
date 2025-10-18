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
	if !value.IsValidCurrency() {
		return false
	}
	return value.IsValidAmount()
}

func (value Money) IsValidAmount() bool {
	return value.Amount >= 0
}

func (value Money) IsValidCurrency() bool {
	switch value.Currency {
	case USD, EUR, JPY:
		return true
	default:
		return false
	}
}

func (value Money) IsLessThan(other Money) bool {
	if value.Currency != other.Currency {
		return false
	}
	return value.Amount < other.Amount
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