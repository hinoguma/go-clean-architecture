package model

import "errors"

type DomainErr struct {
	Type ErrorType
	err  error
}

func (e DomainErr) Error() string {
	return e.err.Error()
}

func (e DomainErr) HasErr() bool {
	return e.err != nil
}

func (e DomainErr) GetErr() error {
	return e.err
}

type ErrorType string

const (
	ErrTypeResourceNotFound ErrorType = "resource_not_found"
	ErrTypeValidation       ErrorType = "validation_error"
	ErrTypeUnauthorized     ErrorType = "unauthorized"
	ErrTypeInternal         ErrorType = "internal_error"
	ErrTypeTimeout          ErrorType = "timeout"
	ErrTypeNotImplemented   ErrorType = "not_implemented"
)

type NotFoundError error

func NewNotFoundError(message string) NotFoundError {
	return NotFoundError(errors.New(message))
}

func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, NotFoundError(nil))
}
