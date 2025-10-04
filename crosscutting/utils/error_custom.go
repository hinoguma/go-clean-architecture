package utils

import "errors"

const (
	ErrorTypeGeneral           ErrorType = "GENERAL"
	ErrorTypeDataNotFound      ErrorType = "DATA_NOT_FOUND"
	ErrorTypeAuthorizeFailed   ErrorType = "AUTHORIZE_FAILED"
	ErrorTypeDataLocked        ErrorType = "DATA_LOCKED"
	ErrorTypeConditionNotMatch ErrorType = "CONDITION_NOT_MATCH"
	ErrorTypeDataConvertFailed ErrorType = "DATA_CONVERT_FAILED"
	ErrorTypeValidationFailed  ErrorType = "VALIDATION_FAILED"
	ErrorTypeUnexpectedFormat  ErrorType = "UNEXPECTED_FORMAT"
	ErrorTypeTimeout           ErrorType = "TIMEOUT"
	ErrorTypeForbidden         ErrorType = "FORBIDDEN"
	ErrorTypePanic             ErrorType = "PANIC"
)

func IsErrorType(err error, et ErrorType) bool {
	target := NewError("")
	target.Type = et
	return errors.Is(err, target)
}

func IsDataNotFoundError(err error) bool {
	return IsErrorType(err, ErrorTypeDataNotFound)
}

func IsDataLockedError(err error) bool {
	return IsErrorType(err, ErrorTypeDataLocked)
}

func IsConditionNotMatchError(err error) bool {
	return IsErrorType(err, ErrorTypeConditionNotMatch)
}

func NewDataNotFoundError(field string, value string) Error {
	err := NewError("data not found")
	err.Type = ErrorTypeDataNotFound
	err.SetAttr("field", field).
		SetAttr("value", value)
	return err
}

func IsTimeoutError(err error) bool {
	return IsErrorType(err, ErrorTypeTimeout)
}

func NewTimeoutError(max string) Error {
	err := NewError("operation timeout")
	err.Type = ErrorTypeTimeout
	err.SetAttr("max", max)
	return err
}

func IsAuthorizeFailedError(err error) bool {
	return IsErrorType(err, ErrorTypeAuthorizeFailed)
}

func NewAuthorizeFailedError(message string) Error {
	err := NewError(message)
	err.Type = ErrorTypeAuthorizeFailed
	return err
}

func NewValidationFailedError(details []ValidateErrorDetail) Error {
	err := NewError("validation failed")
	err.Type = ErrorTypeValidationFailed
	err.SetAttr("details", details)
	return err
}
