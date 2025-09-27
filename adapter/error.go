package adapter

import "app/usecase"

type Error struct {
	Type  ErrorType
	error error
}

type ErrorType string

const (
	ErrTypeResourceNotFound   ErrorType = "resource_not_found"
	ErrTypeValidation         ErrorType = "validation_error"
	ErrTypeUnauthorized       ErrorType = "unauthorized"
	ErrTypeForbidden          ErrorType = "forbidden"
	ErrTypeTimeout            ErrorType = "timeout"
	ErrTypeInternal           ErrorType = "internal_error"
	ErrTypeServiceUnavailable ErrorType = "service_unavailable"
	ErrTypeNotImplemented     ErrorType = "not_implemented"
	ErrTypeUseCaseError       ErrorType = "usecase_error"
)

func ConvertUseCaseToAdapterErr(useCaseErr usecase.Error) Error {
	t, ok := useCaseErrToAdapterErrMap[useCaseErr.Type]
	if !ok {
		t = ErrTypeUseCaseError
	}
	return Error{
		Type:  t,
		error: useCaseErr.GetErr(),
	}
}

var useCaseErrToAdapterErrMap = map[usecase.ErrorType]ErrorType{
	usecase.ErrTypeResourceNotFound:   ErrTypeResourceNotFound,
	usecase.ErrTypeValidation:         ErrTypeValidation,
	usecase.ErrTypeUnauthorized:       ErrTypeUnauthorized,
	usecase.ErrTypeForbidden:          ErrTypeForbidden,
	usecase.ErrTypeTimeout:            ErrTypeTimeout,
	usecase.ErrTypeInternal:           ErrTypeInternal,
	usecase.ErrTypeServiceUnavailable: ErrTypeServiceUnavailable,
	usecase.ErrTypeNotImplemented:     ErrTypeNotImplemented,
}
