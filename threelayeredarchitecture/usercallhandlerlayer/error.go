package usercallhandlerlayer

import "fmt"

type ValidationErrorDetail struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

func (ve ValidationErrorDetail) JsonString() string {
	return fmt.Sprintf(`{"field":"%s","reason":"%s"}`, ve.Field, ve.Reason)
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

func NewRequiredValidationError(field string) ValidationErrorDetail {
	return ValidationErrorDetail{
		Field:  field,
		Reason: fmt.Sprintf(`%s is required`, field),
	}
}

func NewEmptyError(field string) ValidationErrorDetail {
	return ValidationErrorDetail{
		Field:  field,
		Reason: fmt.Sprintf(`%s is empty`, field),
	}
}

func NewInvalidValueError(field string) ValidationErrorDetail {
	return ValidationErrorDetail{
		Field:  field,
		Reason: fmt.Sprintf(`%s is invalid value`, field),
	}
}

func NewMustBeStringError(field string) ValidationErrorDetail {
	return ValidationErrorDetail{
		Field:  field,
		Reason: fmt.Sprintf(`%s must be string`, field),
	}
}

func NewMustBeIntegerError(field string) ValidationErrorDetail {
	return ValidationErrorDetail{
		Field:  field,
		Reason: fmt.Sprintf(`%s must be string`, field),
	}
}

func NewMustBeRangeError(field string, from, to int64) ValidationErrorDetail {
	return ValidationErrorDetail{
		Field:  field,
		Reason: fmt.Sprintf(`%s must be %d to %d`, field, from, to),
	}
}

type ErrorResponseBody struct {
	Message string                  `json:"message"`
	Details []ValidationErrorDetail `json:"details,omitempty"`
}

func (body ErrorResponseBody) JsonString() string {
	if len(body.Details) <= 0 {
		return fmt.Sprintf(`{"message":"%s"}`, body.Message)
	}
	details := ""
	for i, detail := range body.Details {
		if i == 0 {
			details = detail.JsonString()
			continue
		}
		details = details + "," + detail.JsonString()
	}
	return fmt.Sprintf(`{"message":"%s","details":[%s]}`, body.Message, details)
}

func NewBadRequestBody() ErrorResponseBody {
	return ErrorResponseBody{
		Message: "Bad request",
		Details: nil,
	}
}

func NewValidateErrorBody(details []ValidationErrorDetail) ErrorResponseBody {
	return ErrorResponseBody{
		Message: "validation error",
		Details: details,
	}
}
