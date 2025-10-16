package adapter

import (
	"app/experimentarchitecture/applogiclayer/usecase"
	utils2 "app/experimentarchitecture/crosscutting/utils"
)

type ErrorType string

const (
	ErrorTypeValidation      ErrorType = "ValidateError"
	ErrorTypeAuthorization   ErrorType = "AuthorizationError"
	ErrorTypeNotFound        ErrorType = "NotFoundError"
	ErrorTypeInternalServer  ErrorType = "InternalServerError"
	ErrorTypeTooManyRequests ErrorType = "TooManyRequestsError"
	ErrorTypeApplication     ErrorType = "ApplicationError"
)

func NewAdapterErrorValidation(errs []ValidateError) AdapterError {
	details := make([]utils2.ValidateErrorDetail, 0, len(errs))
	for _, ve := range errs {
		details = append(details, ve.UtilValidateErrorDetail())
	}
	return AdapterError{
		ErrType:        ErrorTypeValidation,
		Err:            utils2.NewValidationFailedError(details),
		ValidateErrors: errs,
	}
}

type AdapterError struct {
	ErrType        ErrorType
	Err            error
	ValidateErrors []ValidateError
}

func (err AdapterError) IsValidateError() bool {
	return err.ErrType == ErrorTypeValidation
}

func (err AdapterError) GetValidateErrors() []ValidateError {
	if err.ValidateErrors == nil {
		return make([]ValidateError, 0)
	}
	return err.ValidateErrors
}

func (err AdapterError) Error() string {
	if err.Err == nil {
		return ""
	}
	return err.Err.Error()
}

func NewValidateErrorRequired(field string) ValidateError {
	return ValidateError{
		Field:   field,
		Message: "is required",
	}
}

func NewValidateErrorMustBeString(field string) ValidateError {
	return ValidateError{
		Field:   field,
		Message: "must be string",
	}
}

func NewValidateError(field string, message string) ValidateError {
	return ValidateError{
		Field:   field,
		Message: message,
	}
}

type ValidateError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (ve *ValidateError) SetByUseCase(detail usecase.ValidateError) *ValidateError {
	ve.Field = detail.Field
	ve.Message = detail.Message
	return ve
}

func (ve ValidateError) UtilValidateErrorDetail() utils2.ValidateErrorDetail {
	return utils2.ValidateErrorDetail{
		Field:   ve.Field,
		Message: ve.Message,
	}
}

func convertValidateErrorsFromUseCase(ucves []usecase.ValidateError) []ValidateError {
	ves := make([]ValidateError, 0, len(ucves))
	for _, ucve := range ucves {
		ve := ValidateError{}
		ve.SetByUseCase(ucve)
		ves = append(ves, ve)
	}
	return ves
}

func handleErrFromUseCase(err error) AdapterError {
	if err == nil {
		return AdapterError{}
	}
	ucErr, ok := usecase.AssertUseCaseError(err)
	if !ok {
		return AdapterError{
			ErrType: ErrorTypeApplication,
			Err:     err,
		}
	}

	if ucErr.IsValidateError() {
		return AdapterError{
			ErrType:        ErrorTypeValidation,
			Err:            err,
			ValidateErrors: convertValidateErrorsFromUseCase(ucErr.GetValidateErrors()),
		}
	}

	if ucErr.IsAuthorizationError() {
		return AdapterError{
			ErrType: ErrorTypeAuthorization,
			Err:     err,
		}
	}

	return AdapterError{
		ErrType: ErrorTypeApplication,
		Err:     err,
	}
}
