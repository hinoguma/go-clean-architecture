package main

import (
	"app/cleanarchitecture/frameworksdrivers/web/apigwlambda/handler"
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// lambda + apigateway proxy

func main() {

	lambda.Start(GetHandler)
}

func GetHandler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {

	switch event.RawPath {
	case "/transfer":
		return handler.TransferHandler(ctx, event)
	default:
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 404,
			Body:       "Not Found",
		}, nil
	}

}
