package appinfralayer

import "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

type dbItemDTO interface {
	setByDynamoDBAttrs(attrs map[string]types.AttributeValue) error
}
