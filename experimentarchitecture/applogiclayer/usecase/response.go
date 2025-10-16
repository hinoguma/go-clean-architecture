package usecase

import (
	"app/experimentarchitecture/applogiclayer/domain/model"
)

// ValidateError
// AuthorizationError
// NotFoundError
// InternalServerError
// TooManyRequestsError
// ApplicationError
//

type ResultType string

const (
	ResultTypeSuccessCompleted ResultType = "SUCCESSFULLY_COMPLETED"
	ResultTypeSuccessAccepted  ResultType = "SUCCESSFULLY_ACCEPTED"
	ResultTypeFailed           ResultType = "FAILED"
)

type BaseUseCaseResult struct {
	ResultType ResultType
}

func (res BaseUseCaseResult) IsSuccess() bool {
	return res.IsSuccessCompleted() || res.IsSuccessAccepted()
}

func (res BaseUseCaseResult) IsSuccessCompleted() bool {
	return res.ResultType == ResultTypeSuccessCompleted
}

func (res BaseUseCaseResult) IsSuccessAccepted() bool {
	return res.ResultType == ResultTypeSuccessAccepted
}

type ErrorType string

const (
	ErrorTypeValidation      ErrorType = "ValidateError"
	ErrorTypeAuthorization   ErrorType = "AuthorizationError"
	ErrorTypeNotFound        ErrorType = "NotFoundError"
	ErrorTypeInternalServer  ErrorType = "InternalServerError"
	ErrorTypeTooManyRequests ErrorType = "TooManyRequestsError"
	ErrorTypeApplication     ErrorType = "ApplicationError"
)

func NewAuthorizationError(err error) UseCaseError {
	return UseCaseError{
		errType: ErrorTypeAuthorization,
		err:     err,
	}
}

func AssertUseCaseError(err error) (UseCaseError, bool) {
	if err == nil {
		return UseCaseError{}, false
	}
	ucErr, ok := err.(UseCaseError)
	return ucErr, ok
}

type UseCaseError struct {
	errType        ErrorType
	err            error
	validateErrors []ValidateError
}

func (err UseCaseError) IsValidateError() bool {
	return err.errType == ErrorTypeValidation
}

func (err UseCaseError) IsAuthorizationError() bool {
	return err.errType == ErrorTypeAuthorization
}

func (err UseCaseError) IsNotFoundError() bool {
	return err.errType == ErrorTypeNotFound
}

func (err UseCaseError) GetValidateErrors() []ValidateError {
	if err.validateErrors == nil {
		return make([]ValidateError, 0)
	}
	return err.validateErrors
}

func (err UseCaseError) Error() string {
	if err.err == nil {
		return ""
	}
	return err.err.Error()
}

type ValidateError struct {
	Field   string
	Message string
}

func (ve *ValidateError) SetByDomainModel(detail model.ValidateError) *ValidateError {
	ve.Field = detail.Field
	ve.Message = detail.Message
	return ve
}
