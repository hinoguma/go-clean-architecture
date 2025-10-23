package usercallhandlerlayer

import (
	"github.com/aws/aws-lambda-go/events"
)

func NewValidateErrResponse(err error) events.APIGatewayV2HTTPResponse {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 400,
	}
}

func NewValidateErrResponseWithDetail(details []ValidationErrorDetail) events.APIGatewayV2HTTPResponse {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 400,
		Body:       NewValidateErrorBody(details).JsonString(),
	}
}

func NewCannotMarshalJsonErrResponse() events.APIGatewayV2HTTPResponse {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 400,
		Body:       NewBadRequestBody().JsonString(),
	}
}

func NewInternalServerErrResponse(err error) events.APIGatewayV2HTTPResponse {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 500,
	}
}

func NewAuthErrResponse(err error) events.APIGatewayV2HTTPResponse {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 401,
	}
}
