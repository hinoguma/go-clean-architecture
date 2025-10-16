package appinfralayer

import (
	"app/experimentarchitecture/appinfrainterfacelayer"
	model2 "app/experimentarchitecture/applogiclayer/domain/model"
	"app/experimentarchitecture/crosscutting/timegenerator"
	utils2 "app/experimentarchitecture/crosscutting/utils"
	"context"
	"database/sql"
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

func NewCustomerUserRepository() appinfrainterfacelayer.CustomerUserRepository {
	return customerUserRepository{}
}

func (r customerUserRepository) Get(ctx context.Context, id model2.CustomerUserID) (model2.CustomerUser, error) {
	var exec SqlExecutor = r.sqlClient.db
	tx := GetTxFromContext(ctx)
	if tx != nil {
		exec = tx
	}
	raw := QueryRowContext(exec, ctx, "SELECT id, name, created_at, updated_at FROM customer_users WHERE id = ?", id)
	if raw.Err() != nil {
		return model2.CustomerUser{}, utils2.ErrWrap(raw.Err(), utils2.NewError("failed to query customer user"))
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
		return model2.CustomerUser{}, utils2.ErrWrap(err, utils2.NewError("failed to scan customer user"))
	}
	return dto.toCustomerUserModel(), nil
}

func (r customerUserRepository) UpdateByID(ctx context.Context, id model2.CustomerUserID, updateReq model2.UpdateItemRequest) (model2.CustomerUser, error) {
	// build input
	input, err := modelUpdateRequestToDynamoDBUpdateItemInput(updateReq)
	if err != nil {
		return model2.CustomerUser{}, err
	}
	input.Key = newIdKey(id.String())
	input.TableName = aws.String(tableCustomerUser)

	// execute update
	output, err := r.dynamoDBClient.UpdateItem(ctx, &input)
	if err != nil {
		return model2.CustomerUser{}, err
	}

	// extract result
	dto := CustomerUserDTO{}
	err = attributevalue.UnmarshalMap(output.Attributes, &dto)
	if err != nil {
		return model2.CustomerUser{}, utils2.ErrWrap(err, utils2.NewError("failed to unmarshal customer user"))
	}
	return dto.toCustomerUserModel(), nil
}

type customerUserDynamoDBRepository struct {
	sqlClient      SqlClient
	dynamoDBClient dynamoDBClient
}

func NewCustomerUserDynamoDBRepository() appinfrainterfacelayer.CustomerUserRepository {
	return customerUserDynamoDBRepository{}
}

func (r customerUserDynamoDBRepository) Get(ctx context.Context, id model2.CustomerUserID) (model2.CustomerUser, error) {

	keyCond := model2.NewIdCondition(id)
	dto := CustomerUserDTO{}
	req := getItemRequest{
		tableName:        tableCustomerUser,
		keyCondition:     keyCond,
		isConsistentRead: aws.Bool(true),
	}
	_, err := r.dynamoDBClient.GetItem(ctx, req, &dto)
	if err != nil {
		return model2.CustomerUser{}, err
	}
	return dto.toCustomerUserModel(), nil
}

func (r customerUserDynamoDBRepository) Create(ctx context.Context, item model2.CustomerUser) error {
	// build input
	itemDto := NewCustomerUserDTOByModel(item)
	av, err := attributevalue.MarshalMap(itemDto)
	if err != nil {
		return utils2.ErrWrap(err, utils2.NewError("failed to marshal customer user"))
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

func (r customerUserDynamoDBRepository) Delete(ctx context.Context, id model2.CustomerUserID) error {
	// build input
	builder := modelItemConditionToDynamoDBConditionExpression(
		notLockedCondition(),
	)
	expr, err := expression.NewBuilder().WithCondition(builder).Build()
	if err != nil {
		return utils2.ErrWrap(err, utils2.NewError("failed to build condition expression"))
	}
	input := dynamodb.DeleteItemInput{
		TableName:           aws.String(tableCustomerUser),
		Key:                 newIdKey(id.String()),
		ConditionExpression: expr.Condition(),
	}

	// execute delete
	_, err = r.dynamoDBClient.DeleteItem(ctx, &input)
	if err != nil {
		return utils2.ErrWrap(err, utils2.NewError("failed to delete customer user"))
	}
	return nil
}

func (r customerUserDynamoDBRepository) UpdateByID(ctx context.Context, id model2.CustomerUserID, updateReq model2.UpdateItemRequest) (model2.CustomerUser, error) {
	// build input
	input, err := modelUpdateRequestToDynamoDBUpdateItemInput(updateReq)
	if err != nil {
		return model2.CustomerUser{}, err
	}
	input.Key = newIdKey(id.String())
	input.TableName = aws.String(tableCustomerUser)

	// execute update
	output, err := r.dynamoDBClient.UpdateItem(ctx, &input)
	if err != nil {
		return model2.CustomerUser{}, err
	}

	// extract result
	dto := CustomerUserDTO{}
	err = attributevalue.UnmarshalMap(output.Attributes, &dto)
	if err != nil {
		return model2.CustomerUser{}, utils2.ErrWrap(err, utils2.NewError("failed to unmarshal customer user"))
	}
	return dto.toCustomerUserModel(), nil
}

func (r customerUserDynamoDBRepository) TxOperationGet(ctx context.Context, id model2.CustomerUserID) (model2.CustomerUser, error) {
	txp, ok := getTxProcessFromContext(ctx)
	if !ok {
		txId := GetTransactionIDFromContext(ctx)
		return model2.CustomerUser{}, utils2.NewError("no transaction in context").
			ToPointer().SetAttr("transactionId", txId).ToValue()
	}
	txLogger := getOperationTxLogger{
		transactionProcess: txp,
		operationId:        utils2.UuidV4(),
		tablename:          tableCustomerUser,
	}
	err := txLogger.startOperationLog()
	if err != nil {
		return model2.CustomerUser{}, err
	}
	item, err := r.txQueryGetByID(ctx, txLogger.startOpLog, id)
	if err != nil {
		logErr := txLogger.completeFailedOperationLog()
		if logErr != nil {
			err = utils2.ErrWrap(errors.New(logErr.Error()), err)
		}
		return model2.CustomerUser{}, err
	}
	item, err = r.txQueryLockItem(ctx, txLogger.startOpLog, item)
	if err != nil {
		logErr := txLogger.completeFailedOperationLog()
		if logErr != nil {
			err = utils2.ErrWrap(errors.New(logErr.Error()), err)
		}
		return model2.CustomerUser{}, err
	}
	return item, nil
}

func (r customerUserDynamoDBRepository) TxOperationCreate(
	ctx context.Context, item model2.CustomerUser,
) error {
	// prepare transaction
	txp, ok := getTxProcessFromContext(ctx)
	if !ok {
		txId := GetTransactionIDFromContext(ctx)
		return utils2.NewError("no transaction in context").
			ToPointer().SetAttr("transactionId", txId).ToValue()
	}
	itemDto := NewCustomerUserDTOByModel(item)
	txLogger := createOperationTxLogger{
		transactionProcess: txp,
		operationId:        utils2.UuidV4(),
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
			err = utils2.ErrWrap(errors.New(logErr.Error()), err)
		}
		return err
	}
	err = txLogger.completeSuccessOperationLog()
	if err != nil {
		logErr := txLogger.completeFailedOperationLog()
		if logErr != nil {
			err = utils2.ErrWrap(errors.New(logErr.Error()), err)
		}
		return err
	}
	return nil
}

func (r customerUserDynamoDBRepository) TxOperationUpdateByID(
	ctx context.Context, id model2.CustomerUserID, updateReq model2.UpdateItemRequest,
) (model2.CustomerUser, error) {
	// prepare transaction
	txp, ok := getTxProcessFromContext(ctx)
	if !ok {
		txId := GetTransactionIDFromContext(ctx)
		return model2.CustomerUser{}, utils2.NewError("no transaction in context").
			ToPointer().SetAttr("transactionId", txId).ToValue()
	}
	txLogger := updateOperationTxLogger{
		transactionProcess: txp,
		operationId:        utils2.UuidV4(),
		operationType:      operationTypeUpdate,
		tablename:          tableCustomerUser,
		updateReq:          updateReq,
	}

	// start operation log
	err := txLogger.startOperationLog()
	if err != nil {
		return model2.CustomerUser{}, err
	}

	// get an old item
	beforeItem, err := r.txQueryGetByID(ctx, txLogger.startOpLog, id)
	if err != nil {
		// query error log
		return model2.CustomerUser{}, err
	}
	_, err = r.txQueryLockItem(ctx, txLogger.startOpLog, beforeItem)
	if err != nil {
		logErr := txLogger.completeFailedOperationLog()
		if logErr != nil {
			err = utils2.ErrWrap(errors.New(logErr.Error()), err)
		}
		return model2.CustomerUser{}, err
	}

	item, err := r.txQueryUpdateByID(ctx, txLogger.startOpLog, id, updateReq)
	if err != nil {
		logErr := txLogger.completeFailedOperationLog()
		if logErr != nil {
			err = utils2.ErrWrap(errors.New(logErr.Error()), err)
		}
		return model2.CustomerUser{}, err
	}

	// complete operation log
	err = txLogger.completeSuccessOperationLog()
	if err != nil {
		logErr := txLogger.completeFailedOperationLog()
		if logErr != nil {
			err = utils2.ErrWrap(errors.New(logErr.Error()), err)
		}
		return model2.CustomerUser{}, err
	}
	return item, nil
}

func (r customerUserDynamoDBRepository) TxOperationDelete(
	ctx context.Context, id model2.CustomerUserID,
) error {
	// prepare transaction
	txp, ok := getTxProcessFromContext(ctx)
	if !ok {
		txId := GetTransactionIDFromContext(ctx)
		return utils2.NewError("no transaction in context").
			ToPointer().SetAttr("transactionId", txId).ToValue()
	}
	txLogger := deleteOperationTxLogger{
		transactionProcess: txp,
		operationId:        utils2.UuidV4(),
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
			err = utils2.ErrWrap(errors.New(logErr.Error()), err)
		}
		return err
	}

	// delete item
	err = r.txQueryDelete(ctx, txLogger.startOpLog, id)
	if err != nil {
		logErr := txLogger.completeFailedOperationLog()
		if logErr != nil {
			err = utils2.ErrWrap(errors.New(logErr.Error()), err)
		}
		return err
	}

	// complete operation log
	err = txLogger.completeSuccessOperationLog()
	if err != nil {
		logErr := txLogger.completeFailedOperationLog()
		if logErr != nil {
			err = utils2.ErrWrap(errors.New(logErr.Error()), err)
		}
		return err
	}
	return nil
}

func (r customerUserDynamoDBRepository) txQueryCreate(
	ctx context.Context, startOpLog startOperationLog, item model2.CustomerUser,
) error {
	// prepare transaction
	txp, ok := getTxProcessFromContext(ctx)
	if !ok {
		txId := GetTransactionIDFromContext(ctx)
		return utils2.NewError("no transaction in context").
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
			err = utils2.ErrWrap(errors.New(logErr.Error()), err)
		}
		return err
	}

	// create item
	err = r.Create(ctx, item)
	if err != nil {
		logErr := txLogger.queryResultFailedLog(startLog, err)
		if logErr != nil {
			err = utils2.ErrWrap(errors.New(logErr.Error()), err)
		}
		return err
	}

	// create result log
	keyCond := model2.ItemConditionSingle{
		Operator: model2.ComparisonOperatorEqual,
		Field:    "id",
		Value:    item.ID.String(),
	}
	err = txLogger.queryResultLog(startLog, keyCond)
	if err != nil {
		logErr := txLogger.queryResultFailedLog(startLog, err)
		if logErr != nil {
			err = utils2.ErrWrap(errors.New(logErr.Error()), err)
		}
		return err
	}
	return nil
}

func (r customerUserDynamoDBRepository) txQueryGetByID(
	ctx context.Context, startOpLog startOperationLog, id model2.CustomerUserID,
) (model2.CustomerUser, error) {
	// prepare transaction
	txp, ok := getTxProcessFromContext(ctx)
	if !ok {
		txId := GetTransactionIDFromContext(ctx)
		return model2.CustomerUser{}, utils2.NewError("no transaction in context").
			ToPointer().SetAttr("transactionId", txId).ToValue()
	}
	txLogger := getQueryLogger{
		ctx:                ctx,
		startOpLog:         startOpLog,
		transactionProcess: txp,
		tablename:          tableCustomerUser,
		condition: model2.ItemConditionSingle{
			Operator: model2.ComparisonOperatorEqual,
			Field:    "id",
			Value:    id.String(),
		},
	}

	// start query log
	startLog, err := txLogger.startQueryLog()
	if err != nil {
		logErr := txLogger.queryResultFailedLog(startLog, err)
		if logErr != nil {
			err = utils2.ErrWrap(errors.New(logErr.Error()), err)
		}
		return model2.CustomerUser{}, err
	}

	// execute to get item
	item, err := r.Get(ctx, id)
	if err != nil {
		logErr := txLogger.queryResultFailedLog(startLog, err)
		if logErr != nil {
			err = utils2.ErrWrap(errors.New(logErr.Error()), err)
		}
		return model2.CustomerUser{}, err
	}

	// query result log
	itemDto := NewCustomerUserDTOByModel(item)
	err = txLogger.queryResultLog(startLog, itemDto.JsonString())
	if err != nil {
		logErr := txLogger.queryResultFailedLog(startLog, err)
		if logErr != nil {
			err = utils2.ErrWrap(errors.New(logErr.Error()), err)
		}
		return model2.CustomerUser{}, err
	}
	return item, nil
}

func (r customerUserDynamoDBRepository) txQueryLockItem(
	ctx context.Context, startOpLog startOperationLog, item model2.CustomerUser,
) (model2.CustomerUser, error) {
	updateReq := newUpdateItemRequestForLock(
		startOpLog.TransactionId, timegenerator.NowTs(), "", item.UpdatedAt,
	)
	item, err := r.txQueryUpdateByID(ctx, startOpLog, item.ID, updateReq)
	if err != nil {
		// update error log
		return model2.CustomerUser{}, err
	}
	return item, nil
}

func (r customerUserDynamoDBRepository) txQueryUpdateByID(
	ctx context.Context, startOpLog startOperationLog, id model2.CustomerUserID, updateReq model2.UpdateItemRequest,
) (model2.CustomerUser, error) {
	// prepare transaction
	txp, ok := getTxProcessFromContext(ctx)
	if !ok {
		txId := GetTransactionIDFromContext(ctx)
		return model2.CustomerUser{}, utils2.NewError("no transaction in context").
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
			err = utils2.ErrWrap(errors.New(logErr.Error()), err)
		}
		return model2.CustomerUser{}, err
	}

	// update item
	item, err := r.UpdateByID(ctx, id, updateReq)
	if err != nil {
		logErr := txLogger.queryResultFailedLog(startLog, err)
		if logErr != nil {
			err = utils2.ErrWrap(errors.New(logErr.Error()), err)
		}
		return model2.CustomerUser{}, err
	}

	// update result log
	itemDto := NewCustomerUserDTOByModel(item)
	err = txLogger.queryResultLog(startLog, itemDto.JsonString())
	if err != nil {
		logErr := txLogger.queryResultFailedLog(startLog, err)
		if logErr != nil {
			err = utils2.ErrWrap(errors.New(logErr.Error()), err)
		}
		return model2.CustomerUser{}, err
	}
	return item, nil
}

func (r customerUserDynamoDBRepository) txQueryDelete(
	ctx context.Context, startOpLog startOperationLog, id model2.CustomerUserID,
) error {
	// prepare transaction
	txp, ok := getTxProcessFromContext(ctx)
	if !ok {
		txId := GetTransactionIDFromContext(ctx)
		return utils2.NewError("no transaction in context").
			ToPointer().SetAttr("transactionId", txId).ToValue()
	}
	txLogger := deleteQueryLogger{
		ctx:                ctx,
		startOpLog:         startOpLog,
		transactionProcess: txp,
		tablename:          tableCustomerUser,
		condition: model2.ItemConditionSingle{
			Operator: model2.ComparisonOperatorEqual,
			Field:    "id",
			Value:    id.String(),
		},
	}

	// start query log
	startLog, err := txLogger.startQueryLog()
	if err != nil {
		logErr := txLogger.queryResultFailedLog(startLog, err)
		if logErr != nil {
			err = utils2.ErrWrap(errors.New(logErr.Error()), err)
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
			err = utils2.ErrWrap(errors.New(logErr.Error()), err)
		}
		return err
	}

	// delete result log
	itemDto := NewCustomerUserDTOByModel(item)
	err = txLogger.queryResultLog(startLog, itemDto.JsonString())
	if err != nil {
		logErr := txLogger.queryResultFailedLog(startLog, err)
		if logErr != nil {
			err = utils2.ErrWrap(errors.New(logErr.Error()), err)
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

func (dto CustomerUserDTO) toCustomerUserModel() model2.CustomerUser {
	return model2.CustomerUser{
		ID:        model2.CustomerUserID(dto.ID),
		Name:      dto.Name,
		CreatedAt: utils2.UnixTimestamp(dto.CreatedAt),
		UpdatedAt: utils2.UnixTimestamp(dto.UpdatedAt),
	}
}

func (dto CustomerUserDTO) JsonString() utils2.JsonString {
	return utils2.JsonString(
		fmt.Sprintf(`{"id":"%s","name":"%s","createdAt":%f,"updatedAt":%f}`,
			dto.ID,
			dto.Name,
			dto.CreatedAt,
			dto.UpdatedAt,
		),
	)
}

func (dto *CustomerUserDTO) setByDynamoDBAttrs(attrs map[string]types.AttributeValue) error {
	err := attributevalue.UnmarshalMap(attrs, dto)
	if err != nil {
		return utils2.ErrWrap(err, utils2.NewError("failed to unmarshal map of types.AttributeValue"))
	}
	return nil
}

func (dto *CustomerUserDTO) setBySQLRaw(raw sql.Row) error {
	err := raw.Scan(
		dto.ID,
		dto.Name,
		dto.CreatedAt,
		dto.UpdatedAt,
	)
	if err != nil {
		return utils2.ErrWrap(err, utils2.NewError("failed to raw.Scan()"))
	}
	return nil
}

func NewCustomerUserDTOByModel(m model2.CustomerUser) CustomerUserDTO {
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
		return CustomerUserDTO{}, utils2.ErrWrap(err, utils2.NewError("failed to unmarshal customer user"))
	}
	return dto, nil
}

func newIdKey(id string) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"id": &types.AttributeValueMemberS{Value: id},
	}
}
