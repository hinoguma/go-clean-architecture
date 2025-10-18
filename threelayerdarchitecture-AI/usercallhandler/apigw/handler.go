package main

import (
    "context"
    "encoding/base64"
    "log"
    "app/threelayerdarchitecture/setup"
    "app/threelayerdarchitecture/usercalladapter/transfer"
    "github.com/aws/aws-lambda-go/events"
)

type TransferHandler struct {
    adapter transfer.Adapter
}

func NewTransferHandler(container *setup.Container) TransferHandler {
    return TransferHandler{adapter: container.NewTransferAdapter()}
}

func (h TransferHandler) Handle(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
    body := []byte(event.Body)
    if event.IsBase64Encoded {
        decoded, decodeErr := base64.StdEncoding.DecodeString(event.Body)
        if decodeErr == nil {
            body = decoded
        }
    }

    request := transfer.HTTPRequest{
        Path:          event.RawPath,
        Body:          body,
        Headers:       event.Headers,
        RequestUserID: extractUserID(event),
        CorrelationID: event.RequestContext.RequestID,
    }

    response, err := h.adapter.Execute(ctx, request)
    if err != nil {
        log.Printf("transfer handler error: %v", err)
    }

    return events.APIGatewayV2HTTPResponse{
        StatusCode: response.StatusCode,
        Body:       response.Body,
        Headers:    response.Headers,
    }, nil
}

func extractUserID(event events.APIGatewayV2HTTPRequest) string {
    if event.RequestContext.Authorizer != nil {
        if lambdaMap := event.RequestContext.Authorizer.Lambda; lambdaMap != nil {
            if user, ok := lambdaMap["userID"].(string); ok {
                return user
            }
        }
        if jwt := event.RequestContext.Authorizer.JWT; jwt != nil {
            if sub, ok := jwt.Claims["sub"]; ok && sub != "" {
                return sub
            }
        }
    }
    if user := event.Headers["x-user-id"]; user != "" {
        return user
    }
    if user := event.Headers["X-User-ID"]; user != "" {
        return user
    }
    return ""
}
