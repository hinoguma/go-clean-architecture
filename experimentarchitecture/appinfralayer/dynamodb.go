package appinfralayer

import (
	"app/experimentarchitecture/applogiclayer/domain/model"
	utils2 "app/experimentarchitecture/crosscutting/utils"
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type dynamoDBClient struct {
	client *dynamodb.Client
}

func (client dynamoDBClient) GetItem(ctx context.Context, req getItemRequest, outputItem dbItemDTO) (*dynamodb.GetItemOutput, error) {
	key, err := modelItemConditionToDynamoDBMapAttrs(req.keyCondition)
	if err != nil {
		return nil, err
	}
	input := dynamodb.GetItemInput{
		Key:            key,
		TableName:      aws.String(req.tableName),
		ConsistentRead: req.isConsistentRead,
	}

	// fetch from db...
	output, err := client.client.GetItem(ctx, &input)
	if err != nil {
		nfe := &types.ResourceNotFoundException{}
		if errors.As(err, &nfe) {
			return output, utils2.ErrWrap(
				err, utils2.NewDataNotFoundError().ToPointer().SetAttr("key", key).ToValue(),
			)
		}
		return output, utils2.ErrWrap(err, utils2.NewError("failed to get item"))
	}

	if output.Item == nil {
		return output, utils2.NewDataNotFoundError().ToPointer().SetAttr("key", key).ToValue()
	}

	err = attributevalue.UnmarshalMap(output.Item, outputItem)
	if err != nil {
		return output, utils2.ErrWrap(err, utils2.NewError("failed to unmarshal item"))
	}

	return output, err
}

func (client dynamoDBClient) PutItem(ctx context.Context, input *dynamodb.PutItemInput, optFns ...func(options *dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	output, err := client.client.PutItem(ctx, input, optFns...)
	if err != nil {
		var cce *types.ConditionalCheckFailedException
		if errors.Is(err, cce) {
			return output, utils2.NewError("item already exists").
				ToPointer().SetAttr("item", input.Item).ToValue()
		}
		return output, utils2.ErrWrap(err, utils2.NewError("failed to put item"))
	}
	return output, nil
}

func (client dynamoDBClient) UpdateItem(ctx context.Context, input *dynamodb.UpdateItemInput, optFns ...func(options *dynamodb.Options)) (*dynamodb.UpdateItemOutput, error) {

	panic("not implemented")
}

func (client dynamoDBClient) DeleteItem(ctx context.Context, input *dynamodb.DeleteItemInput, optFns ...func(options *dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {

	output, err := client.client.DeleteItem(ctx, input, optFns...)
	if err != nil {
		var cce *types.ConditionalCheckFailedException
		if errors.As(err, &cce) {
			return output, utils2.NewError("condition not match").
				ToPointer().SetAttr("input", input).ToValue()
		}
	}
	return output, nil
}

func (client dynamoDBClient) Undo(ctx context.Context, undo undoOperation) error {

	if undo.OperationType == operationTypeDelete {
		err := client.deleteByUndo(ctx, undo)
		if err != nil {
			return err
		}
	}

	return nil
}

func (client dynamoDBClient) deleteByUndo(ctx context.Context, undo undoOperation) error {

	if undo.Key == nil {
		return utils2.NewError("undo,Key not exists")
	}
	key, err := modelItemConditionToDynamoDBMapAttrs(*undo.Key)
	if err != nil {
		return err
	}
	input := dynamodb.DeleteItemInput{
		Key:       key,
		TableName: aws.String(undo.TableName),
	}

	if undo.Condition != nil {
		condBuilder := modelItemConditionToDynamoDBConditionExpression(*undo.Condition)
		expr, err := expression.NewBuilder().WithCondition(condBuilder).Build()
		if err != nil {
			return utils2.ErrWrap(err, utils2.NewError("failed to build expression"))
		}
		expr.Condition()
	}

	_, err = client.DeleteItem(ctx, &input)
	if err != nil {
		return err
	}
	return nil

}

func modelUpdateRequestToDynamoDBUpdateItemInput(updateReq model.UpdateItemRequest) (dynamodb.UpdateItemInput, error) {
	input := dynamodb.UpdateItemInput{
		Key:                                 nil,
		TableName:                           nil,
		AttributeUpdates:                    nil,
		ConditionExpression:                 nil,
		ConditionalOperator:                 "",
		Expected:                            nil,
		ExpressionAttributeNames:            nil,
		ExpressionAttributeValues:           nil,
		ReturnConsumedCapacity:              "",
		ReturnItemCollectionMetrics:         "",
		ReturnValues:                        "",
		ReturnValuesOnConditionCheckFailure: "",
		UpdateExpression:                    nil,
	}
	// convert model.UpdateItemRequest to dynamodb.UpdateItemInput
	key, err := modelItemConditionToDynamoDBMapAttrs(updateReq.KeyCondition)
	if err != nil {
		return dynamodb.UpdateItemInput{}, err
	}
	input.Key = key

	// UpdateExpression
	updateBuilder := modelUpdateValueRequestsToDynamoDBUpdateBuilder(updateReq.Values)
	expr, err := expression.NewBuilder().WithUpdate(updateBuilder).Build()
	if err != nil {
		err := utils2.ErrWrap(err, utils2.NewError("failed to build update expression"))
		return input, err
	}
	input.UpdateExpression = expr.Update()
	input.ExpressionAttributeNames = expr.Names()
	input.ExpressionAttributeValues = expr.Values()

	// ConditionalExpression
	if updateReq.HasCondition() {
		condExpr := modelItemConditionToDynamoDBConditionExpression(updateReq.Condition)
		expr, err = expression.NewBuilder().WithCondition(condExpr).Build()
		if err != nil {
			err := utils2.ErrWrap(err, utils2.NewError("failed to build condition expression"))
			return input, err
		}
		input.ConditionExpression = expr.Condition()
	}

	return input, nil
}

const (
	ReadConsistencyLevelStrong   readConsistencyLevel = "STRONG"
	ReadConsistencyLevelEventual readConsistencyLevel = "EVENTUAL"
)

type getItemRequest struct {
	tableName        string
	keyCondition     model.ItemCondition
	outputItem       dbItemDTO
	isConsistentRead *bool
}

type getItemResult struct {
}

func modelItemConditionToDynamoDBMapAttrs(condition model.ItemCondition) (map[string]types.AttributeValue, error) {
	attrs := make(map[string]types.AttributeValue)
	switch condition.(type) {
	case model.ItemConditionSingle:
		c := condition.(model.ItemConditionSingle)
		av, err := attributevalue.Marshal(c.Value)
		if err != nil {
			return nil, utils2.ErrWrap(err, utils2.NewError("failed to marshal condition value"))
		}
		attrs[c.Field] = av
	case model.ItemConditionAnd:
		andCond := condition.(model.ItemConditionAnd)
		for _, cond := range andCond.Conditions() {
			subAttrs, err := modelItemConditionToDynamoDBMapAttrs(cond)
			if err != nil {
				return nil, err
			}
			for k, v := range subAttrs {
				attrs[k] = v
			}
		}
	case model.ItemConditionOr:
		orCond := condition.(model.ItemConditionOr)
		for _, cond := range orCond.Conditions() {
			subAttrs, err := modelItemConditionToDynamoDBMapAttrs(cond)
			if err != nil {
				return nil, err
			}
			for k, v := range subAttrs {
				attrs[k] = v
			}
		}
	}
	return attrs, nil
}

func modelItemConditionToDynamoDBConditionExpression(condition model.ItemCondition) expression.ConditionBuilder {
	builder := expression.ConditionBuilder{}
	switch condition.(type) {
	case model.ItemConditionSingle:
		c := condition.(model.ItemConditionSingle)
		switch c.Operator {
		case model.ComparisonOperatorEqual:
			return expression.Name(c.Field).Equal(expression.Value(c.Value))
		default:
			return expression.Name(c.Field).Equal(expression.Value(c.Value))
		}
	case model.ItemConditionAnd:
		andCond := condition.(model.ItemConditionAnd)
		for _, cond := range andCond.Conditions() {
			builder = builder.And(modelItemConditionToDynamoDBConditionExpression(cond))
		}
		return builder
	case model.ItemConditionOr:
		orCond := condition.(model.ItemConditionOr)
		for _, cond := range orCond.Conditions() {
			builder = builder.Or(modelItemConditionToDynamoDBConditionExpression(cond))
		}
		return builder
	}
	return builder
}

func modelUpdateValueRequestsToDynamoDBUpdateBuilder(reqs []model.UpdateValueRequest) expression.UpdateBuilder {
	updateBuilder := expression.UpdateBuilder{}
	for _, req := range reqs {
		if req.IsSet() {
			updateBuilder.Set(expression.Name(req.Field), expression.Value(req.Value))
		}
		if req.IsAdd() {
			updateBuilder.Add(expression.Name(req.Field), expression.Value(req.Value))
		}
		if req.IsRemove() {
			updateBuilder.Remove(expression.Name(req.Field))
		}
		if req.IsDelete() {
			updateBuilder.Delete(expression.Name(req.Field), expression.Value(req.Value))
		}
	}
	return updateBuilder
}

func notLockedCondition() model.ItemCondition {
	return model.ItemConditionSingle{
		Operator: model.ComparisonOperatorEqual,
		Field:    "txId",
		Value:    "",
	}
}

const attrNotExistsID = "attribute_not_exists(id)"
