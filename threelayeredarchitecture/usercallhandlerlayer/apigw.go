package usercallhandlerlayer

import "github.com/aws/aws-lambda-go/events"

func NewValidateErrResponseUnderAdapter() events.APIGatewayV2HTTPResponse {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: StatusCodeValidateError,
		Body:       NewBadRequestBody().JsonString(),
	}
}

func NewValidateErrResponseWithDetail(details []ValidationErrorDetail) events.APIGatewayV2HTTPResponse {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: StatusCodeValidateError,
		Body:       NewValidateErrorBody(details).JsonString(),
	}
}

func NewCannotMarshalJsonErrResponse() events.APIGatewayV2HTTPResponse {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: StatusCodeValidateError,
		Body:       NewCannotParseRequestBody().JsonString(),
	}
}

func NewInternalServerErrResponse(err error) events.APIGatewayV2HTTPResponse {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: StatusCodeInternalAppError,
	}
}

func NewAuthErrResponse(err error) events.APIGatewayV2HTTPResponse {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: StatusCodeUnAuth,
	}
}
