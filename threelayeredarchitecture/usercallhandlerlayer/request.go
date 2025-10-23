package usercallhandlerlayer

type Currency string

func (value Currency) IsValid() bool {
	switch value {
	case "USD", "EUR", "JPY":
		return true
	}
	return false
}
