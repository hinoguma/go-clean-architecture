package model

type AuthErrorReason struct {
	TokenExpired  bool
	InvalidToken  bool
	InternalError bool
}

type BankCustomerAuthRequest struct {
	Token string
}

type BankCustomerAuthResult struct {
	Customer    BankCustomer
	ErrorReason AuthErrorReason
	Err         error
}

func (res BankCustomerAuthResult) IsSuccess() bool {
	return res.Err == nil
}

func (res BankCustomerAuthResult) IsAuthenticateFailedError() bool {
	return res.ErrorReason.TokenExpired || res.ErrorReason.InvalidToken
}
