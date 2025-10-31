package apigateway

import (
	"app/threelayeredarchitecture/setuplayer/registry"
	"context"
	"github.com/aws/aws-lambda-go/events"
)

func transferHandler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	ucHandler := registry.RegistryUserCallHandler.GetTransferHandler()
	return ucHandler.Execute(ctx, event)
}
