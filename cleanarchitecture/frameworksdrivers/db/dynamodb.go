package db

import (
	"app/cleanarchitecture/crosscutting/errors"
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoDBClient interface {
	GetItem(ctx context.Context, input dynamodb.GetItemInput) (dynamodb.GetItemOutput, error)
	Query(ctx context.Context, input dynamodb.QueryInput) (dynamodb.QueryOutput, error)
	PutItem(ctx context.Context, input dynamodb.PutItemInput, optFns ...func(options *dynamodb.Options)) (dynamodb.PutItemOutput, error)
	UpdateItem(ctx context.Context, input dynamodb.UpdateItemInput, optFns ...func(options *dynamodb.Options)) (dynamodb.UpdateItemOutput, error)
	DeleteItem(ctx context.Context, input dynamodb.DeleteItemInput, optFns ...func(options *dynamodb.Options)) (dynamodb.DeleteItemOutput, error)
}

func NewDynamoDBClient(ctx context.Context) (DynamoDBClient, error) {
	cnf, err := NewAWSConfig(ctx)
	if err != nil {
		return nil, errors.LiftWithCtx(err, ctx)
	}
	return &dynamoDBClient{
		dynamo: dynamodb.NewFromConfig(cnf),
	}, nil
}

type dynamoDBClient struct {
	dynamo *dynamodb.Client
}

func (d *dynamoDBClient) GetItem(ctx context.Context, input dynamodb.GetItemInput) (dynamodb.GetItemOutput, error) {
	output, err := d.dynamo.GetItem(ctx, &input)
	if err != nil {
		return dynamodb.GetItemOutput{}, errors.LiftWithCtx(err, ctx)
	}
	if output == nil {
		return dynamodb.GetItemOutput{}, errors.NewWithCtx("output is empty", ctx)
	}
	return *output, nil
}

func (d *dynamoDBClient) Query(ctx context.Context, input dynamodb.QueryInput) (dynamodb.QueryOutput, error) {
	output, err := d.dynamo.Query(ctx, &input)
	if err != nil {
		return dynamodb.QueryOutput{}, errors.LiftWithCtx(err, ctx)
	}
	if output == nil {
		return dynamodb.QueryOutput{}, errors.NewWithCtx("output is empty", ctx)
	}
	return *output, nil
}

func (d *dynamoDBClient) PutItem(ctx context.Context, input dynamodb.PutItemInput, optFns ...func(options *dynamodb.Options)) (dynamodb.PutItemOutput, error) {
	output, err := d.dynamo.PutItem(ctx, &input, optFns...)
	if err != nil {
		return dynamodb.PutItemOutput{}, errors.LiftWithCtx(err, ctx)
	}
	if output == nil {
		return dynamodb.PutItemOutput{}, errors.NewWithCtx("output is empty", ctx)
	}
	return *output, nil
}

func (d *dynamoDBClient) UpdateItem(ctx context.Context, input dynamodb.UpdateItemInput, optFns ...func(options *dynamodb.Options)) (dynamodb.UpdateItemOutput, error) {
	output, err := d.dynamo.UpdateItem(ctx, &input, optFns...)
	if err != nil {
		return dynamodb.UpdateItemOutput{}, errors.LiftWithCtx(err, ctx)
	}
	if output == nil {
		return dynamodb.UpdateItemOutput{}, errors.NewWithCtx("output is empty", ctx)
	}
	return *output, nil
}

func (d *dynamoDBClient) DeleteItem(ctx context.Context, input dynamodb.DeleteItemInput, optFns ...func(options *dynamodb.Options)) (dynamodb.DeleteItemOutput, error) {
	output, err := d.dynamo.DeleteItem(ctx, &input, optFns...)
	if err != nil {
		return dynamodb.DeleteItemOutput{}, errors.LiftWithCtx(err, ctx)
	}
	if output == nil {
		return dynamodb.DeleteItemOutput{}, errors.NewWithCtx("output is empty", ctx)
	}
	return *output, nil
}
