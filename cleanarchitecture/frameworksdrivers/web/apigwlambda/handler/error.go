package handler

import "github.com/aws/aws-lambda-go/events"

func errorResponse(err error) (events.APIGatewayV2HTTPResponse, error) {

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 500,
		Body:       "Internal Server Error",
	}, nil
}
