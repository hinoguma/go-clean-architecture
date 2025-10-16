package appinfralayer

import (
	"app/experimentarchitecture/appinfrainterfacelayer"
	"app/experimentarchitecture/applogiclayer/domain/model"
	"app/experimentarchitecture/crosscutting/logger"
	"app/experimentarchitecture/crosscutting/timegenerator"
	utils2 "app/experimentarchitecture/crosscutting/utils"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
)

var txList transactionList = transactionList{processMap: make(map[string]transactionProcess)}

type transactionProcess struct {
	transactionId string
	tx            *sql.Tx
}

type transactionList struct {
	processMap map[string]transactionProcess
}

func (tl transactionList) add(tr transactionProcess) error {
	if tl.processMap == nil {
		tl.processMap = make(map[string]transactionProcess)
	}
	_, ok := tl.processMap[tr.transactionId]
	if ok {
		err := utils2.NewError("transaction already exists")
		err.SetAttr("transactionId", tr.transactionId)
		return err
	}
	tl.processMap[tr.transactionId] = tr
	return nil
}

func (tl transactionList) get(transactionId string) (transactionProcess, bool) {
	if tl.processMap == nil {
		return transactionProcess{}, false
	}
	tr, ok := tl.processMap[transactionId]
	return tr, ok
}

func (tl transactionList) remove(transactionId string) {
	if tl.processMap == nil {
		return
	}
	delete(tl.processMap, transactionId)
}

func NewSQLTransactionFactory() appinfrainterfacelayer.TransactionFactory {
	return sqlTransactionFactory{}
}

type sqlTransactionFactory struct {
	// e.g db client...
}

func (factory sqlTransactionFactory) BeginTx(ctx context.Context, opt model.TxOptions) (context.Context, appinfrainterfacelayer.Transaction, error) {
	client, err := NewSqlClient(sqlConfig{})
	if err != nil {
		return ctx, nil, err
	}
	tx, err := client.BeginTx(ctx, opt)
	if err != nil {
		return ctx, nil, err
	}
	err = txList.add(
		transactionProcess{
			transactionId: opt.TransactionID,
			tx:            tx,
		},
	)
	if err != nil {
		return ctx, nil, err
	}
	//
	ctx = SetTransactionIDToContext(ctx, opt.TransactionID)
	return ctx, transaction{transactionId: opt.TransactionID}, nil
}

type transaction struct {
	transactionId string
}

func (t transaction) GetTransactionID() string {
	return t.transactionId
}

func (t transaction) Commit() error {
	tr, ok := txList.get(t.transactionId)
	if !ok {
		return NewTransactionNotFoundError(t.transactionId)
	}
	err := tr.tx.Commit()
	if err != nil {
		return err
	}
	txList.remove(t.transactionId)
	return nil
}

func (t transaction) Rollback() error {
	tr, ok := txList.get(t.transactionId)
	if !ok {
		return NewTransactionNotFoundError(t.transactionId)
	}
	err := tr.tx.Rollback()
	if err != nil {
		return err
	}
	txList.remove(t.transactionId)
	return nil
}

func NewTransactionNotFoundError(transactionId string) error {
	err := utils2.NewError(fmt.Sprintf("transaction not found"))
	err.SetAttr("transactionId", transactionId)
	return err
}

func SetTransactionIDToContext(ctx context.Context, transactionId string) context.Context {
	if ctx == nil {
		return nil
	}
	return context.WithValue(ctx, "transactionId", transactionId)
}

func GetTransactionIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	txId, ok := ctx.Value("transactionId").(string)
	if !ok {
		return ""
	}
	return txId
}

func GetTxFromContext(ctx context.Context) *sql.Tx {
	if ctx == nil {
		return nil
	}
	transactionId := GetTransactionIDFromContext(ctx)
	if transactionId == "" {
		return nil
	}
	tr, ok := txList.get(transactionId)
	if !ok {
		return nil
	}
	return tr.tx
}

func getTxProcessFromContext(ctx context.Context) (transactionProcess, bool) {
	if ctx == nil {
		return transactionProcess{}, false
	}
	transactionId := GetTransactionIDFromContext(ctx)
	if transactionId == "" {
		return transactionProcess{}, false
	}
	tr, ok := txList.get(transactionId)
	return tr, ok
}

type operationType string

const (
	operationTypeInsert operationType = "INSERT"
	operationTypePut    operationType = "PUT"
	operationTypeUpdate operationType = "UPDATE"
	operationTypeDelete operationType = "DELETE"
	operationTypeSelect operationType = "SELECT"
)

type operationErrorType string

const (
	operationErrorTypeNotFoundData      operationErrorType = "NOT_FOUND_DATA"
	operationErrorTypeDataLocked        operationErrorType = "DATA_LOCKED"
	operationErrorTypeConditionNotMatch operationErrorType = "CONDITION_NOT_MATCH"
	operationErrorTypeError             operationErrorType = "ERROR"
)

type transactionLogType string

const (
	transactionLogTypeStartOperation    transactionLogType = "START_OPERATION"
	transactionLogTypeCompleteOperation transactionLogType = "COMPLETE_OPERATION"
	transactionLogTypeStartQuery        transactionLogType = "START_QUERY"
	transactionLogTypeQueryResult       transactionLogType = "QUERY_RESULT"
)

type transactionLog interface {
	LogType() transactionLogType
}

type commonPropsOfLog struct {
	TransactionId string             `json:"transactionId"`
	OperationId   string             `json:"operationId"`
	Type          transactionLogType `json:"logType"`
	Time          utils2.Ymd_Hms_ms  `json:"timestamp"`
}

func (v commonPropsOfLog) LogType() transactionLogType {
	return v.Type
}

type startOperationLog struct {
	commonPropsOfLog
	OperationType operationType `json:"operationType"`
	TableName     string        `json:"tableName"`
}

func (log startOperationLog) JsonString() (utils2.JsonString, error) {
	b, err := json.Marshal(log)
	if err != nil {
		return "{}", utils2.ErrWrap(err, utils2.NewError("failed to marshal startOperationLog"))
	}
	return utils2.JsonString(b), nil
}

func newStartOperationLog(txId string, opId string, opType operationType, tn string, t utils2.Ymd_Hms_ms) startOperationLog {
	return startOperationLog{
		commonPropsOfLog: commonPropsOfLog{
			TransactionId: txId,
			OperationId:   opId,
			Type:          transactionLogTypeStartOperation,
			Time:          t,
		},
		OperationType: opType,
		TableName:     tn,
	}
}

type completeOperationLog struct {
	commonPropsOfLog
	IsSuccess     bool          `json:"isSuccess"`
	OperationType operationType `json:"operationType"`
}

func newCompleteOperationLog(startOpLog startOperationLog, isSuccess bool, t utils2.Ymd_Hms_ms) completeOperationLog {
	return completeOperationLog{
		commonPropsOfLog: commonPropsOfLog{
			TransactionId: startOpLog.TransactionId,
			OperationId:   startOpLog.OperationId,
			Type:          startOpLog.Type,
			Time:          t,
		},
		OperationType: startOpLog.OperationType,
		IsSuccess:     isSuccess,
	}
}

func (v completeOperationLog) JsonString() (utils2.JsonString, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "{}", utils2.ErrWrap(err, utils2.NewError("failed to marshal completeOperationLog"))
	}
	return utils2.JsonString(b), nil
}

type commonPropsOfQueryLog struct {
	QueryId   string        `json:"queryId"`
	QueryType operationType `json:"queryType"`
	TableName string        `json:"tableName"`
}

// todo: need to be independent on domain model?
type startQueryLog struct {
	commonPropsOfLog
	commonPropsOfQueryLog
	Condition         *model.ItemCondition     `json:"condition,omitempty"`
	UpdateItemRequest *model.UpdateItemRequest `json:"updateItemRequest,omitempty"`
	ItemData          utils2.JsonString        `json:"ItemData,omitempty"`
}

func newStartQueryLog(
	startOpLog startOperationLog, qId string, qType operationType, tn string, t utils2.Ymd_Hms_ms) *startQueryLog {
	return &startQueryLog{
		commonPropsOfLog: commonPropsOfLog{
			TransactionId: startOpLog.TransactionId,
			OperationId:   startOpLog.OperationId,
			Type:          transactionLogTypeStartQuery,
			Time:          t,
		},
		commonPropsOfQueryLog: commonPropsOfQueryLog{
			QueryId:   qId,
			QueryType: qType,
			TableName: tn,
		},
	}
}

func (v *startQueryLog) WithCondition(cond model.ItemCondition) *startQueryLog {
	v.Condition = &cond
	return v
}

func (v *startQueryLog) WithUpdateItemRequest(req model.UpdateItemRequest) *startQueryLog {
	v.UpdateItemRequest = &req
	return v
}

func (v *startQueryLog) WithItemData(data utils2.JsonString) *startQueryLog {
	v.ItemData = data
	return v
}

func (v startQueryLog) JsonString() (utils2.JsonString, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "{}", utils2.ErrWrap(err, utils2.NewError("failed to marshal queryResultLog"))
	}
	return utils2.JsonString(b), nil
}

func (v *startQueryLog) ToValue() startQueryLog {
	return *v
}

type queryResultLog struct {
	commonPropsOfLog
	commonPropsOfQueryLog
	IsSuccess     bool                `json:"isSuccess"`
	ErrorType     *operationErrorType `json:"errorType,omitempty"`
	ItemData      *utils2.JsonString  `json:"itemData,omitempty"`
	UndoOperation *undoOperation      `json:"undoOperation,omitempty"`
}

func newQueryResultLog(
	startQueryLog startQueryLog, isSuccess bool, t utils2.Ymd_Hms_ms,
) *queryResultLog {
	return &queryResultLog{
		commonPropsOfLog: commonPropsOfLog{
			TransactionId: startQueryLog.TransactionId,
			OperationId:   startQueryLog.OperationId,
			Type:          transactionLogTypeQueryResult,
			Time:          t,
		},
		commonPropsOfQueryLog: commonPropsOfQueryLog{
			QueryId:   startQueryLog.QueryId,
			QueryType: startQueryLog.QueryType,
			TableName: startQueryLog.TableName,
		},
		IsSuccess: isSuccess,
	}
}

func (v *queryResultLog) WithIsSuccess(isSuccess bool) *queryResultLog {
	v.IsSuccess = isSuccess
	return v
}

func (v *queryResultLog) WithErrorType(errType operationErrorType) *queryResultLog {
	v.ErrorType = &errType
	return v
}

func (v *queryResultLog) WithItemData(data utils2.JsonString) *queryResultLog {
	v.ItemData = &data
	return v
}

func (v *queryResultLog) WithUndoOperation(op undoOperation) *queryResultLog {
	v.UndoOperation = &op
	return v
}

func (v queryResultLog) JsonString() (utils2.JsonString, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "{}", utils2.ErrWrap(err, utils2.NewError("failed to marshal queryResultLog"))
	}
	return utils2.JsonString(b), nil
}

func (v *queryResultLog) ToValue() queryResultLog {
	return *v
}

type undoOperation struct {
	OperationType     operationType            `json:"operationType"`
	TableName         string                   `json:"tableName"`
	Key               *model.ItemCondition     `json:"key,omitempty"`
	ItemData          *utils2.JsonString       `json:"itemData,omitempty"`
	UpdateItemRequest *model.UpdateItemRequest `json:"updateItemRequest,omitempty"`
	Condition         *model.ItemCondition     `json:"condition,omitempty"`
}

func newUndoOperation(opType operationType, tn string) *undoOperation {
	return &undoOperation{
		OperationType: opType,
		TableName:     tn,
	}
}

func (v *undoOperation) WithKey(cond model.ItemCondition) *undoOperation {
	v.Key = &cond
	return v
}

func (v *undoOperation) WithItemData(data utils2.JsonString) *undoOperation {
	v.ItemData = &data
	return v
}

func (v *undoOperation) WithUpdateItemRequest(req model.UpdateItemRequest) *undoOperation {
	v.UpdateItemRequest = &req
	return v
}

func (v *undoOperation) WithCondition(cond model.ItemCondition) *undoOperation {
	v.Condition = &cond
	return v
}

func (v undoOperation) JsonString() (utils2.JsonString, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "{}", utils2.ErrWrap(err, utils2.NewError("failed to marshal queryResultLog"))
	}
	return utils2.JsonString(b), nil
}

func (v *undoOperation) ToValue() undoOperation {
	return *v
}

var fileLogger = logger.NewStdLogger()

type updateOperationTxLogger struct {
	ctx                context.Context
	transactionProcess transactionProcess
	operationId        string
	operationType      operationType
	tablename          string
	updateReq          model.UpdateItemRequest
	startOpLog         startOperationLog
}

func (txLogger updateOperationTxLogger) startOperationLog() error {
	startOpLog := newStartOperationLog(
		txLogger.transactionProcess.transactionId,
		txLogger.operationId,
		txLogger.operationType,
		txLogger.tablename,
		timegenerator.NowYmd_Hms_ms(),
	)
	message, err := startOpLog.JsonString()
	if err != nil {
		return err
	}
	fileLogger.Info(txLogger.ctx, utils2.NewLogRequest(message.String()).ToValue())
	return nil
}

func (txLogger updateOperationTxLogger) completeSuccessOperationLog() error {
	completeOpLog := newCompleteOperationLog(
		txLogger.startOpLog,
		true,
		timegenerator.NowYmd_Hms_ms(),
	)
	message, err := completeOpLog.JsonString()
	if err != nil {
		return err
	}
	fileLogger.Info(txLogger.ctx, utils2.NewLogRequest(message.String()).ToValue())
	return nil
}

func (txLogger updateOperationTxLogger) completeFailedOperationLog() error {
	completeOpLog := newCompleteOperationLog(
		txLogger.startOpLog,
		false,
		timegenerator.NowYmd_Hms_ms(),
	)
	message, err := completeOpLog.JsonString()
	if err != nil {
		return err
	}
	fileLogger.Info(txLogger.ctx, utils2.NewLogRequest(message.String()).ToValue())
	return nil
}

type createOperationTxLogger struct {
	ctx                context.Context
	transactionProcess transactionProcess
	operationId        string
	tablename          string
	item               utils2.JsonString
	startOpLog         startOperationLog
}

func (txLogger createOperationTxLogger) startOperationLog() error {
	startOpLog := newStartOperationLog(
		txLogger.transactionProcess.transactionId,
		txLogger.operationId,
		operationTypeInsert,
		txLogger.tablename,
		timegenerator.NowYmd_Hms_ms(),
	)
	message, err := startOpLog.JsonString()
	if err != nil {
		return err
	}
	fileLogger.Info(txLogger.ctx, utils2.NewLogRequest(message.String()).ToValue())
	return nil
}

func (txLogger createOperationTxLogger) completeSuccessOperationLog() error {
	completeOpLog := newCompleteOperationLog(
		txLogger.startOpLog,
		true,
		timegenerator.NowYmd_Hms_ms(),
	)
	message, err := completeOpLog.JsonString()
	if err != nil {
		return err
	}
	fileLogger.Info(txLogger.ctx, utils2.NewLogRequest(message.String()).ToValue())
	return nil
}

func (txLogger createOperationTxLogger) completeFailedOperationLog() error {
	completeOpLog := newCompleteOperationLog(
		txLogger.startOpLog,
		false,
		timegenerator.NowYmd_Hms_ms(),
	)
	message, err := completeOpLog.JsonString()
	if err != nil {
		return err
	}
	fileLogger.Info(txLogger.ctx, utils2.NewLogRequest(message.String()).ToValue())
	return nil
}

type deleteOperationTxLogger struct {
	ctx                context.Context
	transactionProcess transactionProcess
	operationId        string
	tablename          string
	startOpLog         startOperationLog
}

func (txLogger deleteOperationTxLogger) startOperationLog() error {
	startOpLog := newStartOperationLog(
		txLogger.transactionProcess.transactionId,
		txLogger.operationId,
		operationTypeDelete,
		txLogger.tablename,
		timegenerator.NowYmd_Hms_ms(),
	)
	message, err := startOpLog.JsonString()
	if err != nil {
		return err
	}
	fileLogger.Info(txLogger.ctx, utils2.NewLogRequest(message.String()).ToValue())
	return nil
}

func (txLogger deleteOperationTxLogger) completeSuccessOperationLog() error {
	completeOpLog := newCompleteOperationLog(
		txLogger.startOpLog,
		true,
		timegenerator.NowYmd_Hms_ms(),
	)
	message, err := completeOpLog.JsonString()
	if err != nil {
		return err
	}
	fileLogger.Info(txLogger.ctx, utils2.NewLogRequest(message.String()).ToValue())
	return nil
}

func (txLogger deleteOperationTxLogger) completeFailedOperationLog() error {
	completeOpLog := newCompleteOperationLog(
		txLogger.startOpLog,
		false,
		timegenerator.NowYmd_Hms_ms(),
	)
	message, err := completeOpLog.JsonString()
	if err != nil {
		return err
	}
	fileLogger.Info(txLogger.ctx, utils2.NewLogRequest(message.String()).ToValue())
	return nil
}

type getOperationTxLogger struct {
	ctx                context.Context
	transactionProcess transactionProcess
	operationId        string
	tablename          string
	startOpLog         startOperationLog
}

func (txLogger getOperationTxLogger) startOperationLog() error {
	startOpLog := newStartOperationLog(
		txLogger.transactionProcess.transactionId,
		txLogger.operationId,
		operationTypeSelect,
		txLogger.tablename,
		timegenerator.NowYmd_Hms_ms(),
	)
	message, err := startOpLog.JsonString()
	if err != nil {
		return err
	}
	fileLogger.Info(txLogger.ctx, utils2.NewLogRequest(message.String()).ToValue())
	return nil
}

func (txLogger getOperationTxLogger) completeSuccessOperationLog() error {
	completeOpLog := newCompleteOperationLog(
		txLogger.startOpLog,
		true,
		timegenerator.NowYmd_Hms_ms(),
	)
	message, err := completeOpLog.JsonString()
	if err != nil {
		return err
	}
	fileLogger.Info(txLogger.ctx, utils2.NewLogRequest(message.String()).ToValue())
	return nil
}

func (txLogger getOperationTxLogger) completeFailedOperationLog() error {
	completeOpLog := newCompleteOperationLog(
		txLogger.startOpLog,
		false,
		timegenerator.NowYmd_Hms_ms(),
	)
	message, err := completeOpLog.JsonString()
	if err != nil {
		return err
	}
	fileLogger.Info(txLogger.ctx, utils2.NewLogRequest(message.String()).ToValue())
	return nil
}

type updateQueryLogger struct {
	ctx                context.Context
	startOpLog         startOperationLog
	transactionProcess transactionProcess
	tablename          string
	updateReq          model.UpdateItemRequest
}

func (txLogger updateQueryLogger) startQueryLog() (startQueryLog, error) {
	queryLog := newStartQueryLog(
		txLogger.startOpLog,
		utils2.UuidV4(),
		operationTypeUpdate,
		txLogger.tablename,
		timegenerator.NowYmd_Hms_ms(),
	).WithUpdateItemRequest(txLogger.updateReq)
	message, err := queryLog.ToValue().JsonString()
	if err != nil {
		return queryLog.ToValue(), err
	}
	fileLogger.Info(txLogger.ctx, utils2.NewLogRequest(message.String()).ToValue())
	return queryLog.ToValue(), nil
}

func (txLogger updateQueryLogger) queryResultLog(
	startQueryLog startQueryLog,
	after utils2.JsonString,
) error {
	undo := newUndoOperation(operationTypePut, txLogger.tablename).
		WithUpdateItemRequest(txLogger.updateReq).
		ToValue()
	updateResultLog := newQueryResultLog(
		startQueryLog, true, timegenerator.NowYmd_Hms_ms(),
	).WithItemData(after).
		WithUndoOperation(undo)
	message, err := updateResultLog.ToValue().JsonString()
	if err != nil {
		return err
	}
	fileLogger.Info(txLogger.ctx, utils2.NewLogRequest(message.String()).ToValue())
	return nil
}

func (txLogger updateQueryLogger) queryResultFailedLog(
	startLog startQueryLog, err error,
) error {
	return queryResultFailedLog(txLogger.ctx, startLog, err)
}

type createQueryLogger struct {
	ctx                context.Context
	startOpLog         startOperationLog
	transactionProcess transactionProcess
	tablename          string
	updateReq          model.UpdateItemRequest
}

func (txLogger createQueryLogger) startQueryLog() (startQueryLog, error) {
	queryLog := newStartQueryLog(
		txLogger.startOpLog,
		utils2.UuidV4(),
		operationTypeInsert,
		txLogger.tablename,
		timegenerator.NowYmd_Hms_ms(),
	).WithUpdateItemRequest(txLogger.updateReq)
	message, err := queryLog.ToValue().JsonString()
	if err != nil {
		return queryLog.ToValue(), err
	}
	fileLogger.Info(txLogger.ctx, utils2.NewLogRequest(message.String()).ToValue())
	return queryLog.ToValue(), nil
}

func (txLogger createQueryLogger) queryResultLog(
	startQueryLog startQueryLog,
	condition model.ItemCondition,
) error {

	undo := newUndoOperation(operationTypeDelete, txLogger.tablename).
		WithCondition(condition).
		ToValue()
	updateResultLog := newQueryResultLog(
		startQueryLog, true, timegenerator.NowYmd_Hms_ms(),
	).WithUndoOperation(undo)
	message, err := updateResultLog.ToValue().JsonString()
	if err != nil {
		return err
	}
	fileLogger.Info(txLogger.ctx, utils2.NewLogRequest(message.String()).ToValue())
	return nil
}

func (txLogger createQueryLogger) queryResultFailedLog(
	startLog startQueryLog, err error,
) error {
	return queryResultFailedLog(txLogger.ctx, startLog, err)
}

func queryResultFailedLog(
	ctx context.Context, startLog startQueryLog, err error,
) error {
	log := newQueryResultLog(
		startLog, false, timegenerator.NowYmd_Hms_ms(),
	).WithErrorType(operationErrorTypeError)
	if utils2.IsDataNotFoundError(err) {
		log = log.WithErrorType(operationErrorTypeNotFoundData)
	} else if utils2.IsDataLockedError(err) {
		log = log.WithErrorType(operationErrorTypeDataLocked)
	} else if utils2.IsConditionNotMatchError(err) {
		log = log.WithErrorType(operationErrorTypeConditionNotMatch)
	}
	message, err := log.ToValue().JsonString()
	if err != nil {
		return err
	}
	fileLogger.Info(ctx, utils2.NewLogRequest(message.String()).ToValue())
	return nil
}

type getQueryLogger struct {
	ctx                context.Context
	startOpLog         startOperationLog
	transactionProcess transactionProcess
	tablename          string
	condition          model.ItemCondition
}

func (txLogger getQueryLogger) startQueryLog() (startQueryLog, error) {
	queryLog := newStartQueryLog(
		txLogger.startOpLog,
		utils2.UuidV4(),
		operationTypeSelect,
		txLogger.tablename,
		timegenerator.NowYmd_Hms_ms(),
	).WithCondition(txLogger.condition)
	message, err := queryLog.ToValue().JsonString()
	if err != nil {
		return queryLog.ToValue(), err
	}
	fileLogger.Info(txLogger.ctx, utils2.NewLogRequest(message.String()).ToValue())
	return queryLog.ToValue(), nil
}

func (txLogger getQueryLogger) queryResultLog(
	startLog startQueryLog, item utils2.JsonString,
) error {
	log := newQueryResultLog(
		startLog, true, timegenerator.NowYmd_Hms_ms(),
	).WithItemData(item)
	message, err := log.ToValue().JsonString()
	if err != nil {
		return err
	}
	fileLogger.Info(txLogger.ctx, utils2.NewLogRequest(message.String()).ToValue())
	return nil

}

func (txLogger getQueryLogger) queryResultFailedLog(
	startLog startQueryLog, err error,
) error {
	return queryResultFailedLog(txLogger.ctx, startLog, err)
}

type deleteQueryLogger struct {
	ctx                context.Context
	startOpLog         startOperationLog
	transactionProcess transactionProcess
	tablename          string
	condition          model.ItemCondition
}

func (txLogger deleteQueryLogger) startQueryLog() (startQueryLog, error) {
	queryLog := newStartQueryLog(
		txLogger.startOpLog,
		utils2.UuidV4(),
		operationTypeDelete,
		txLogger.tablename,
		timegenerator.NowYmd_Hms_ms(),
	).WithCondition(txLogger.condition)
	message, err := queryLog.ToValue().JsonString()
	if err != nil {
		return queryLog.ToValue(), err
	}
	fileLogger.Info(txLogger.ctx, utils2.NewLogRequest(message.String()).ToValue())
	return queryLog.ToValue(), nil
}

func (txLogger deleteQueryLogger) queryResultLog(
	startLog startQueryLog, item utils2.JsonString,
) error {
	undo := newUndoOperation(operationTypePut, txLogger.tablename).
		WithItemData(item).
		ToValue()
	log := newQueryResultLog(
		startLog, true, timegenerator.NowYmd_Hms_ms(),
	).WithUndoOperation(undo)
	message, err := log.ToValue().JsonString()
	if err != nil {
		return err
	}
	fileLogger.Info(txLogger.ctx, utils2.NewLogRequest(message.String()).ToValue())
	return nil

}

func (txLogger deleteQueryLogger) queryResultFailedLog(
	startLog startQueryLog, err error,
) error {
	return queryResultFailedLog(txLogger.ctx, startLog, err)
}

func newUpdateItemRequestForLock(
	txId string, ts utils2.UnixTimestamp, isolationLevel string, updatedAt utils2.UnixTimestamp,
) model.UpdateItemRequest {
	return model.UpdateItemRequest{
		Values: []model.UpdateValueRequest{
			{
				Field: "txId",
				Value: txId,
			},
			{
				Field: "txStartedAt",
				Value: ts.Int64(),
			},
			{
				Field: "isolationLevel",
				Value: isolationLevel,
			},
		},
		Condition: model.ItemConditionSingle{
			Operator: model.ComparisonOperatorEqual,
			Field:    "updatedAt",
			Value:    updatedAt,
		},
	}
}

/**
Select
	pattern
			found
				lock record
					update to set tx id -> undo
				no lock
					-> undo update to remove tx id
			not found -> no undo
	details
		start log
			table name
			condition
		execute to select from database
		select result log
			raw data
		update to set tx id
			-> execute UPDATE with select operation id
		complete operation log


INSERT
	pattern
		-> undo delete
	details
		start log
			table name
			insert data
		execute to insert into database
		insert result log
			table name
			inserted data
			undo operation
				type delete
				key condition
		complete operation log

UPDATE
	pattern
		-> undo update with old values
	details
		start log
			operation id
			table name
			update value requests
			condition
		start select log
			table name
			query id
			condition
		execute to select old values from database without lock
		select result log
			query id
			raw data
		start update log
			query id
			table name
			update value requests
			condition
		execute to update into database
		update result log
			query id
			table name
			updated data
			condition
			undo operation
				table name
				type update
				key condition
				update value requests with old values
		complete operation log

PUT
	pattern
		found -> undo put with old values
		not found -> undo delete
	details
		start operation log
			operation id
			table name
		start select log
			query id
			table name
			condition
		execute to select old values from database without lock
		select result log
			query id
			raw data
		start put log
			query id
			table name
			put data
		execute to put into database
		put result log
			query id
			table name
			put data
			undo operation
				if found
					type put
					key condition
					put data with old values
				if not found
					type delete
					key condition
		complete operation log


DELETE
	pattern
		found -> undo put with old values
		not found -> no undo
	details
		start operation log
			operation id
			table name
			condition
		start select log
			query id
			table name
			condition
		execute to select old values from database without lock
		select result log
			query id
			raw data
		start delete log
			query id
			table name
			condition
		execute to delete from database
		delete result log
			query id
			table name
			condition
			undo operation
				if found
					type put
					key condition
					put data with old values
				if not found
					no undo
		complete operation log

*/
