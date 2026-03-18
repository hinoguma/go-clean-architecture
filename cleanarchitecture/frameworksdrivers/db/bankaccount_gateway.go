package db

import (
	"app/cleanarchitecture/crosscutting/errors"
	"app/cleanarchitecture/enterprisebusinessrules/entities"
	"app/cleanarchitecture/interfaceadapters/gateways"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const TableBankAccounts string = "bank_accounts"

type bankAccountGateway struct {
	dynamoDBClient DynamoDBClient
}

func NewBankAccountGateway() gateways.BankAccountGateway {
	return &bankAccountGateway{}
}

func (gateway bankAccountGateway) Get(ctx context.Context, id string) (entities.BankAccount, error) {
	input := dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
		TableName: aws.String(TableBankAccounts),
	}
	res, err := gateway.dynamoDBClient.GetItem(ctx, input)
	if err != nil {
		return entities.BankAccount{}, errors.LiftWithCtx(err, ctx)
	}
	dto := bankAccountDynamoDTO{}
	err = dto.SetByDynamoDBAttrs(res.Item)
	if err != nil {
		return entities.BankAccount{}, errors.LiftWithCtx(err, ctx)
	}
	return dto.ToEntity(), nil
}

func (gateway bankAccountGateway) Create(ctx context.Context, item entities.BankAccount) (string, error) {
	dto := bankAccountDynamoDTO{}
	err := dto.SetByEntity(item)
	if err != nil {
		return "", errors.LiftWithCtx(err, ctx)
	}

	attrs, err := dto.ToDynamoDBAttrs()
	if err != nil {
		return "", errors.LiftWithCtx(err, ctx)
	}

	input := dynamodb.PutItemInput{
		Item:      attrs,
		TableName: aws.String(TableBankAccounts),
	}
	_, err = gateway.dynamoDBClient.PutItem(ctx, input)
	if err != nil {
		return "", errors.LiftWithCtx(err, ctx)
	}
	return dto.ID, nil
}

func (gateway bankAccountGateway) Update(ctx context.Context, item entities.BankAccount) error {
	//TODO implement me
	panic("implement me")
}

func (gateway bankAccountGateway) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

type bankAccountDynamoDTO struct {
	ID              string `dynamodbav:"id"`
	UserID          string `dynamodbav:"user_id"`
	BalanceAmount   int64  `dynamodbav:"balance_amount"`
	BalanceCurrency string `dynamodbav:"balance_currency"`
	CreatedAt       int64  `dynamodbav:"created_at,unixtime"`
	UpdatedAt       int64  `dynamodbav:"updated_at,unixtime"`
}

func (item *bankAccountDynamoDTO) SetByDynamoDBAttrs(attrs map[string]types.AttributeValue) error {
	if attrs == nil {
		return nil
	}
	return attributevalue.UnmarshalMap(attrs, item)
}

func (item *bankAccountDynamoDTO) SetByEntity(entity entities.BankAccount) error {
	item.ID = entity.ID
	item.UserID = entity.UserID
	item.BalanceAmount = entity.Balance.Amount
	item.BalanceCurrency = entity.Balance.Currency
	item.CreatedAt = entity.CreatedAt
	item.UpdatedAt = entity.UpdatedAt
	return nil
}

func (item bankAccountDynamoDTO) ToDynamoDBAttrs() (map[string]types.AttributeValue, error) {
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return nil, errors.LiftWithCtx(err, context.Background())
	}
	return av, nil
}

func (item bankAccountDynamoDTO) ToEntity() entities.BankAccount {
	return entities.BankAccount{
		ID:     item.ID,
		UserID: item.UserID,
		Balance: entities.Money{
			Amount:   item.BalanceAmount,
			Currency: item.BalanceCurrency,
		},
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}
