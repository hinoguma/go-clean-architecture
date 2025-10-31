package model

import (
	"app/threelayeredarchitecture/crosscuttinglayer"
	"errors"
	"fmt"
)

// todo: think about crosscuttinglayer errors and if need, fix this implementation

type ErrorCode string

const (
	ErrorCodeUnauthorized ErrorCode = "UNAUTHORIZED_ERROR"
	ErrorCodeValidation   ErrorCode = "VALIDATION_ERROR"
	ErrorCodeInternal     ErrorCode = "INTERNAL_ERROR"
)

type AppLogicError struct {
	Code             ErrorCode
	Message          string
	Err              error
	ValidationErrors []ValidationErrorDetail
}

func (e *AppLogicError) Error() string {
	return fmt.Sprintf(`{"code":"%s","message":"%s","error":"%s","validation_errors":%v}`, e.Code, e.Message, e.Err, e.ValidationErrors)
}

func NewAppLogicError(
	err error,
) *AppLogicError {
	return &AppLogicError{
		Code: ErrorCodeInternal,
		Err:  err,
	}
}

func (e *AppLogicError) WithCode(
	code ErrorCode,
) *AppLogicError {
	e.Code = code
	return e
}

func (e *AppLogicError) WithMessage(
	message string,
) *AppLogicError {
	e.Message = message
	return e
}

func (e *AppLogicError) WithValidationErrorDetails(
	validationErrors []ValidationErrorDetail,
) *AppLogicError {
	e.ValidationErrors = validationErrors
	return e
}

func NewValidationError(
	details []ValidationErrorDetail,
) *AppLogicError {
	return &AppLogicError{
		Code:             ErrorCodeValidation,
		ValidationErrors: details,
		Err:              errors.New("validation error"),
	}
}

type ValidationErrorDetail struct {
	Field  string
	Reason string
}

func NewValidationErrorDetail(
	field string,
	reason string,
) ValidationErrorDetail {
	return ValidationErrorDetail{
		Field:  field,
		Reason: reason,
	}
}

func NewRequiredErrDetail(
	field string,
) ValidationErrorDetail {
	return ValidationErrorDetail{
		Field:  field,
		Reason: "required",
	}
}

func NewInvalidDataErrDetail(
	field string,
) ValidationErrorDetail {
	return ValidationErrorDetail{
		Field:  field,
		Reason: "invalid_data",
	}
}

func NewMinValueErrDetail(
	field string,
	minValue int,
) ValidationErrorDetail {
	return ValidationErrorDetail{
		Field:  field,
		Reason: fmt.Sprintf("%s has to be greater than %d", field, minValue),
	}
}

func IsNotFoundErrInInfraAdapter(err error) bool {
	return crosscuttinglayer.IsDataNotFound(err)
}
