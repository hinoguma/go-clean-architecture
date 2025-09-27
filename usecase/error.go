package usecase

import "app/domain/model"

type HasError struct {
	Error Error
}

type Error struct {
	Type ErrorType
	err  error
}

func (e Error) Error() string {
	if e.HasErr() == false {
		return ""
	}
	return e.err.Error()
}

func (e Error) HasErr() bool {
	return e.err != nil
}

func (e Error) GetErr() error {
	return e.err
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
	ErrTypeDomainError        ErrorType = "domain_error"
)

func ConvertDomainToUseCaseErr(domainErr model.DomainErr) Error {
	errType, ok := domainErrToUseCaseErrMap[domainErr.Type]
	if !ok {
		errType = ErrTypeDomainError
	}
	return Error{
		Type: errType,
		err:  domainErr.GetErr(),
	}
}

var domainErrToUseCaseErrMap = map[model.ErrorType]ErrorType{
	model.ErrTypeResourceNotFound: ErrTypeResourceNotFound,
	model.ErrTypeValidation:       ErrTypeValidation,
	model.ErrTypeUnauthorized:     ErrTypeUnauthorized,
	model.ErrTypeInternal:         ErrTypeInternal,
	model.ErrTypeTimeout:          ErrTypeTimeout,
	model.ErrTypeNotImplemented:   ErrTypeNotImplemented,
}
