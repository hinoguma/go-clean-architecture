package appinfralayer

import (
	"app/applogiclayer/domain/infrainterface"
	"app/applogiclayer/domain/model"
	"app/crosscutting/timegenerator"
	"app/crosscutting/utils"
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type customerUserRepository struct {
	sqlClient      SqlClient
	dynamoDBClient dynamoDBClient
}

func NewCustomerUserRepository() infrainterface.CustomerUserRepository {
	return customerUserRepository{}
}

func (r customerUserRepository) Get(ctx context.Context, id model.CustomerUserID) (model.CustomerUser, error) {
	var exec SqlExecutor = r.sqlClient.db
	tx := GetTxFromContext(ctx)
	if tx != nil {
		exec = tx
	}
	raw := QueryRowContext(exec, ctx, "SELECT id, name, created_at, updated_at FROM customer_users WHERE id = ?", id)
	if raw.Err() != nil {
		return model.CustomerUser{}, utils.ErrWrap(raw.Err(), utils.NewError("failed to query customer user"))
	}
	// for demonstration, return dummy data
	dto := CustomerUserDTO{}
	err := raw.Scan(
		&dto.ID,
		&dto.Name,
		&dto.CreatedAt,
		&dto.UpdatedAt,
	)
	if err != nil {
		return model.CustomerUser{}, utils.ErrWrap(err, utils.NewError("failed to scan customer user"))
	}
	return dto.toCustomerUserModel(), nil
}

func (r customerUserRepository) UpdateByID(ctx context.Context, id model.CustomerUserID, updateReq model.UpdateItemRequest) (model.CustomerUser, error) {
	// build input
	input, err := modelUpdateRequestToDynamoDBUpdateItemInput(updateReq)
	if err != nil {
		return model.CustomerUser{}, err
	}
	input.Key = newIdKey(id.String())
	input.TableName = aws.String(tableCustomerUser)

	// execute update
	output, err := r.dynamoDBClient.UpdateItem(ctx, &input)
	if err != nil {
		return model.CustomerUser{}, err
	}

	// extract result
	dto := CustomerUserDTO{}
	err = attributevalue.UnmarshalMap(output.Attributes, &dto)
	if err != nil {
		return model.CustomerUser{}, utils.ErrWrap(err, utils.NewError("failed to unmarshal customer user"))
	}
	return dto.toCustomerUserModel(), nil
}

type customerUserDynamoDBRepository struct {
	sqlClient      SqlClient
	dynamoDBClient dynamoDBClient
}

func NewCustomerUserDynamoDBRepository() infrainterface.CustomerUserRepository {
	return customerUserDynamoDBRepository{}
}

func (r customerUserDynamoDBRepository) Get(ctx context.Context, id model.CustomerUserID) (model.CustomerUser, error) {
	var exec SqlExecutor = r.sqlClient.db
	tx := GetTxFromContext(ctx)
	if tx != nil {
		exec = tx
	}
	raw := QueryRowContext(exec, ctx, "SELECT id, name, created_at, updated_at FROM customer_users WHERE id = ?", id)
	if raw.Err() != nil {
		return model.CustomerUser{}, utils.ErrWrap(raw.Err(), utils.NewError("failed to query customer user"))
	}
	// for demonstration, return dummy data
	dto := CustomerUserDTO{}
	err := raw.Scan(
		&dto.ID,
		&dto.Name,
		&dto.CreatedAt,
		&dto.UpdatedAt,
	)
	if err != nil {
		return model.CustomerUser{}, utils.ErrWrap(err, utils.NewError("failed to scan customer user"))
	}
	return dto.toCustomerUserModel(), nil
}

func (r customerUserDynamoDBRepository) Create(ctx context.Context, item model.CustomerUser) error {
	// build input
	itemDto := NewCustomerUserDTOByModel(item)
	av, err := attributevalue.MarshalMap(itemDto)
	if err != nil {
		return utils.ErrWrap(err, utils.NewError("failed to marshal customer user"))
	}
	input := dynamodb.PutItemInput{
		TableName:           aws.String(tableCustomerUser),
		Item:                av,
		ConditionExpression: aws.String(attrNotExistsID),
	}
	// execute put
	_, err = r.dynamoDBClient.PutItem(ctx, &input)
	if err != nil {
		return err
	}
	return nil
}

func (r customerUserDynamoDBRepository) Delete(ctx context.Context, id model.CustomerUserID) error {
	// build input
	builder := modelItemConditionToDynamoDBConditionExpression(
		notLockedCondition(),
	)
	expr, err := expression.NewBuilder().WithCondition(builder).Build()
	if err != nil {
		return utils.ErrWrap(err, utils.NewError("failed to build condition expression"))
	}
	input := dynamodb.DeleteItemInput{
		TableName:           aws.String(tableCustomerUser),
		Key:                 newIdKey(id.String()),
		ConditionExpression: expr.Condition(),
	}

	// execute delete
	_, err = r.dynamoDBClient.DeleteItem(ctx, &input)
	if err != nil {
		return utils.ErrWrap(err, utils.NewError("failed to delete customer user"))
	}
	return nil
}

func (r customerUserDynamoDBRepository) UpdateByID(ctx context.Context, id model.CustomerUserID, updateReq model.UpdateItemRequest) (model.CustomerUser, error) {
	// build input
	input, err := modelUpdateRequestToDynamoDBUpdateItemInput(updateReq)
	if err != nil {
		return model.CustomerUser{}, err
	}
	input.Key = newIdKey(id.String())
	input.TableName = aws.String(tableCustomerUser)

	// execute update
	output, err := r.dynamoDBClient.UpdateItem(ctx, &input)
	if err != nil {
		return model.CustomerUser{}, err
	}

	// extract result
	dto := CustomerUserDTO{}
	err = attributevalue.UnmarshalMap(output.Attributes, &dto)
	if err != nil {
		return model.CustomerUser{}, utils.ErrWrap(err, utils.NewError("failed to unmarshal customer user"))
	}
	return dto.toCustomerUserModel(), nil
}

func (r customerUserDynamoDBRepository) TxOperationGet(ctx context.Context, id model.CustomerUserID) (model.CustomerUser, error) {
	txp, ok := getTxProcessFromContext(ctx)
	if !ok {
		txId := GetTransactionIDFromContext(ctx)
		return model.CustomerUser{}, utils.NewError("no transaction in context").
			ToPointer().SetAttr("transactionId", txId).ToValue()
	}
	txLogger := getOperationTxLogger{
		transactionProcess: txp,
		operationId:        utils.UuidV4(),
		tablename:          tableCustomerUser,
	}
	err := txLogger.startOperationLog()
	if err != nil {
		return model.CustomerUser{}, err
	}
	item, err := r.txQueryGetByID(ctx, txLogger.startOpLog, id)
	if err != nil {
		logErr := txLogger.completeFailedOperationLog()
		if logErr != nil {
			err = utils.ErrWrap(errors.New(logErr.Error()), err)
		}
		return model.CustomerUser{}, err
	}
	item, err = r.txQueryLockItem(ctx, txLogger.startOpLog, item)
	if err != nil {
		logErr := txLogger.completeFailedOperationLog()
		if logErr != nil {
			err = utils.ErrWrap(errors.New(logErr.Error()), err)
		}
		return model.CustomerUser{}, err
	}
	return item, nil
}

func (r customerUserDynamoDBRepository) TxOperationCreate(
	ctx context.Context, item model.CustomerUser,
) error {
	// prepare transaction
	txp, ok := getTxProcessFromContext(ctx)
	if !ok {
		txId := GetTransactionIDFromContext(ctx)
		return utils.NewError("no transaction in context").
			ToPointer().SetAttr("transactionId", txId).ToValue()
	}
	itemDto := NewCustomerUserDTOByModel(item)
	txLogger := createOperationTxLogger{
		transactionProcess: txp,
		operationId:        utils.UuidV4(),
		tablename:          tableCustomerUser,
		item:               itemDto.JsonString(),
	}
	err := txLogger.startOperationLog()
	if err != nil {
		return err
	}
	err = r.txQueryCreate(ctx, txLogger.startOpLog, item)
	if err != nil {
		logErr := txLogger.completeFailedOperationLog()
		if logErr != nil {
			err = utils.ErrWrap(errors.New(logErr.Error()), err)
		}
		return err
	}
	err = txLogger.completeSuccessOperationLog()
	if err != nil {
		logErr := txLogger.completeFailedOperationLog()
		if logErr != nil {
			err = utils.ErrWrap(errors.New(logErr.Error()), err)
		}
		return err
	}
	return nil
}

func (r customerUserDynamoDBRepository) TxOperationUpdateByID(
	ctx context.Context, id model.CustomerUserID, updateReq model.UpdateItemRequest,
) (model.CustomerUser, error) {
	// prepare transaction
	txp, ok := getTxProcessFromContext(ctx)
	if !ok {
		txId := GetTransactionIDFromContext(ctx)
		return model.CustomerUser{}, utils.NewError("no transaction in context").
			ToPointer().SetAttr("transactionId", txId).ToValue()
	}
	txLogger := updateOperationTxLogger{
		transactionProcess: txp,
		operationId:        utils.UuidV4(),
		operationType:      operationTypeUpdate,
		tablename:          tableCustomerUser,
		updateReq:          updateReq,
	}

	// start operation log
	err := txLogger.startOperationLog()
	if err != nil {
		return model.CustomerUser{}, err
	}

	// get an old item
	beforeItem, err := r.txQueryGetByID(ctx, txLogger.startOpLog, id)
	if err != nil {
		// query error log
		return model.CustomerUser{}, err
	}
	_, err = r.txQueryLockItem(ctx, txLogger.startOpLog, beforeItem)
	if err != nil {
		logErr := txLogger.completeFailedOperationLog()
		if logErr != nil {
			err = utils.ErrWrap(errors.New(logErr.Error()), err)
		}
		return model.CustomerUser{}, err
	}

	item, err := r.txQueryUpdateByID(ctx, txLogger.startOpLog, id, updateReq)
	if err != nil {
		logErr := txLogger.completeFailedOperationLog()
		if logErr != nil {
			err = utils.ErrWrap(errors.New(logErr.Error()), err)
		}
		return model.CustomerUser{}, err
	}

	// complete operation log
	err = txLogger.completeSuccessOperationLog()
	if err != nil {
		logErr := txLogger.completeFailedOperationLog()
		if logErr != nil {
			err = utils.ErrWrap(errors.New(logErr.Error()), err)
		}
		return model.CustomerUser{}, err
	}
	return item, nil
}

func (r customerUserDynamoDBRepository) TxOperationDelete(
	ctx context.Context, id model.CustomerUserID,
) error {
	// prepare transaction
	txp, ok := getTxProcessFromContext(ctx)
	if !ok {
		txId := GetTransactionIDFromContext(ctx)
		return utils.NewError("no transaction in context").
			ToPointer().SetAttr("transactionId", txId).ToValue()
	}
	txLogger := deleteOperationTxLogger{
		transactionProcess: txp,
		operationId:        utils.UuidV4(),
		tablename:          tableCustomerUser,
	}

	// start operation log
	err := txLogger.startOperationLog()
	if err != nil {
		return err
	}

	// get an old item
	beforeItem, err := r.txQueryGetByID(ctx, txLogger.startOpLog, id)
	if err != nil {
		// query error log
		return err
	}
	_, err = r.txQueryLockItem(ctx, txLogger.startOpLog, beforeItem)
	if err != nil {
		logErr := txLogger.completeFailedOperationLog()
		if logErr != nil {
			err = utils.ErrWrap(errors.New(logErr.Error()), err)
		}
		return err
	}

	// delete item
	err = r.txQueryDelete(ctx, txLogger.startOpLog, id)
	if err != nil {
		logErr := txLogger.completeFailedOperationLog()
		if logErr != nil {
			err = utils.ErrWrap(errors.New(logErr.Error()), err)
		}
		return err
	}

	// complete operation log
	err = txLogger.completeSuccessOperationLog()
	if err != nil {
		logErr := txLogger.completeFailedOperationLog()
		if logErr != nil {
			err = utils.ErrWrap(errors.New(logErr.Error()), err)
		}
		return err
	}
	return nil
}

func (r customerUserDynamoDBRepository) txQueryCreate(
	ctx context.Context, startOpLog startOperationLog, item model.CustomerUser,
) error {
	// prepare transaction
	txp, ok := getTxProcessFromContext(ctx)
	if !ok {
		txId := GetTransactionIDFromContext(ctx)
		return utils.NewError("no transaction in context").
			ToPointer().SetAttr("transactionId", txId).ToValue()
	}
	txLogger := createQueryLogger{
		ctx:                ctx,
		startOpLog:         startOpLog,
		transactionProcess: txp,
		tablename:          tableCustomerUser,
	}

	// start query log
	startLog, err := txLogger.startQueryLog()
	if err != nil {
		logErr := txLogger.queryResultFailedLog(startLog, err)
		if logErr != nil {
			err = utils.ErrWrap(errors.New(logErr.Error()), err)
		}
		return err
	}

	// create item
	err = r.Create(ctx, item)
	if err != nil {
		logErr := txLogger.queryResultFailedLog(startLog, err)
		if logErr != nil {
			err = utils.ErrWrap(errors.New(logErr.Error()), err)
		}
		return err
	}

	// create result log
	keyCond := model.ItemConditionSingle{
		Operator: model.ComparisonOperatorEqual,
		Field:    "id",
		Value:    item.ID.String(),
	}
	err = txLogger.queryResultLog(startLog, keyCond)
	if err != nil {
		logErr := txLogger.queryResultFailedLog(startLog, err)
		if logErr != nil {
			err = utils.ErrWrap(errors.New(logErr.Error()), err)
		}
		return err
	}
	return nil
}

func (r customerUserDynamoDBRepository) txQueryGetByID(
	ctx context.Context, startOpLog startOperationLog, id model.CustomerUserID,
) (model.CustomerUser, error) {
	// prepare transaction
	txp, ok := getTxProcessFromContext(ctx)
	if !ok {
		txId := GetTransactionIDFromContext(ctx)
		return model.CustomerUser{}, utils.NewError("no transaction in context").
			ToPointer().SetAttr("transactionId", txId).ToValue()
	}
	txLogger := getQueryLogger{
		ctx:                ctx,
		startOpLog:         startOpLog,
		transactionProcess: txp,
		tablename:          tableCustomerUser,
		condition: model.ItemConditionSingle{
			Operator: model.ComparisonOperatorEqual,
			Field:    "id",
			Value:    id.String(),
		},
	}

	// start query log
	startLog, err := txLogger.startQueryLog()
	if err != nil {
		logErr := txLogger.queryResultFailedLog(startLog, err)
		if logErr != nil {
			err = utils.ErrWrap(errors.New(logErr.Error()), err)
		}
		return model.CustomerUser{}, err
	}

	// execute to get item
	item, err := r.Get(ctx, id)
	if err != nil {
		logErr := txLogger.queryResultFailedLog(startLog, err)
		if logErr != nil {
			err = utils.ErrWrap(errors.New(logErr.Error()), err)
		}
		return model.CustomerUser{}, err
	}

	// query result log
	itemDto := NewCustomerUserDTOByModel(item)
	err = txLogger.queryResultLog(startLog, itemDto.JsonString())
	if err != nil {
		logErr := txLogger.queryResultFailedLog(startLog, err)
		if logErr != nil {
			err = utils.ErrWrap(errors.New(logErr.Error()), err)
		}
		return model.CustomerUser{}, err
	}
	return item, nil
}

func (r customerUserDynamoDBRepository) txQueryLockItem(
	ctx context.Context, startOpLog startOperationLog, item model.CustomerUser,
) (model.CustomerUser, error) {
	updateReq := newUpdateItemRequestForLock(
		startOpLog.TransactionId, timegenerator.NowTs(), "", item.UpdatedAt,
	)
	item, err := r.txQueryUpdateByID(ctx, startOpLog, item.ID, updateReq)
	if err != nil {
		// update error log
		return model.CustomerUser{}, err
	}
	return item, nil
}

func (r customerUserDynamoDBRepository) txQueryUpdateByID(
	ctx context.Context, startOpLog startOperationLog, id model.CustomerUserID, updateReq model.UpdateItemRequest,
) (model.CustomerUser, error) {
	// prepare transaction
	txp, ok := getTxProcessFromContext(ctx)
	if !ok {
		txId := GetTransactionIDFromContext(ctx)
		return model.CustomerUser{}, utils.NewError("no transaction in context").
			ToPointer().SetAttr("transactionId", txId).ToValue()
	}
	// start query log
	txLogger := updateQueryLogger{
		ctx:                ctx,
		transactionProcess: txp,
		startOpLog:         startOpLog,
		tablename:          tableCustomerUser,
		updateReq:          updateReq,
	}
	startLog, err := txLogger.startQueryLog()
	if err != nil {
		logErr := txLogger.queryResultFailedLog(startLog, err)
		if logErr != nil {
			err = utils.ErrWrap(errors.New(logErr.Error()), err)
		}
		return model.CustomerUser{}, err
	}

	// update item
	item, err := r.UpdateByID(ctx, id, updateReq)
	if err != nil {
		logErr := txLogger.queryResultFailedLog(startLog, err)
		if logErr != nil {
			err = utils.ErrWrap(errors.New(logErr.Error()), err)
		}
		return model.CustomerUser{}, err
	}

	// update result log
	itemDto := NewCustomerUserDTOByModel(item)
	err = txLogger.queryResultLog(startLog, itemDto.JsonString())
	if err != nil {
		logErr := txLogger.queryResultFailedLog(startLog, err)
		if logErr != nil {
			err = utils.ErrWrap(errors.New(logErr.Error()), err)
		}
		return model.CustomerUser{}, err
	}
	return item, nil
}

func (r customerUserDynamoDBRepository) txQueryDelete(
	ctx context.Context, startOpLog startOperationLog, id model.CustomerUserID,
) error {
	// prepare transaction
	txp, ok := getTxProcessFromContext(ctx)
	if !ok {
		txId := GetTransactionIDFromContext(ctx)
		return utils.NewError("no transaction in context").
			ToPointer().SetAttr("transactionId", txId).ToValue()
	}
	txLogger := deleteQueryLogger{
		ctx:                ctx,
		startOpLog:         startOpLog,
		transactionProcess: txp,
		tablename:          tableCustomerUser,
		condition: model.ItemConditionSingle{
			Operator: model.ComparisonOperatorEqual,
			Field:    "id",
			Value:    id.String(),
		},
	}

	// start query log
	startLog, err := txLogger.startQueryLog()
	if err != nil {
		logErr := txLogger.queryResultFailedLog(startLog, err)
		if logErr != nil {
			err = utils.ErrWrap(errors.New(logErr.Error()), err)
		}
		return err
	}

	// get item for rollback
	item, err := r.txQueryGetByID(ctx, startOpLog, id)

	// delete item
	err = r.Delete(ctx, id)
	if err != nil {
		logErr := txLogger.queryResultFailedLog(startLog, err)
		if logErr != nil {
			err = utils.ErrWrap(errors.New(logErr.Error()), err)
		}
		return err
	}

	// delete result log
	itemDto := NewCustomerUserDTOByModel(item)
	err = txLogger.queryResultLog(startLog, itemDto.JsonString())
	if err != nil {
		logErr := txLogger.queryResultFailedLog(startLog, err)
		if logErr != nil {
			err = utils.ErrWrap(errors.New(logErr.Error()), err)
		}
		return err
	}
	return nil
}

type CustomerUserDTO struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt int64  `json:"updatedAt"`
}

func (dto CustomerUserDTO) toCustomerUserModel() model.CustomerUser {
	return model.CustomerUser{
		ID:        model.CustomerUserID(dto.ID),
		Name:      dto.Name,
		CreatedAt: utils.UnixTimestamp(dto.CreatedAt),
		UpdatedAt: utils.UnixTimestamp(dto.UpdatedAt),
	}
}

func (dto CustomerUserDTO) JsonString() utils.JsonString {
	return utils.JsonString(
		fmt.Sprintf(`{"id":"%s","name":"%s","createdAt":%f,"updatedAt":%f}`,
			dto.ID,
			dto.Name,
			dto.CreatedAt,
			dto.UpdatedAt,
		),
	)
}

func NewCustomerUserDTOByModel(m model.CustomerUser) CustomerUserDTO {
	return CustomerUserDTO{
		ID:        string(m.ID),
		Name:      m.Name,
		CreatedAt: m.CreatedAt.Int64(),
		UpdatedAt: m.UpdatedAt.Int64(),
	}
}

func NewCustomerUserDTOByDynamoAttrs(attrs map[string]types.AttributeValue) (CustomerUserDTO, error) {
	dto := CustomerUserDTO{}
	err := attributevalue.UnmarshalMap(attrs, &dto)
	if err != nil {
		return CustomerUserDTO{}, utils.ErrWrap(err, utils.NewError("failed to unmarshal customer user"))
	}
	return dto, nil
}

func newIdKey(id string) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"id": &types.AttributeValueMemberS{Value: id},
	}
}
