package appinfralayer

import (
	"app/threelayeredarchitecture/config"
	"app/threelayeredarchitecture/crosscuttinglayer"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type BankCustomerRepositoryIF interface {
	Get(ctx context.Context, id string) (BankCustomerRawData, error)
	Create(ctx context.Context, item BankCustomerRawData) error
	Update(ctx context.Context, id string, updateFields UpdateFieldRequests) error
	Delete(ctx context.Context, id string) error

	// With Transaction
	Lock(ctx context.Context, id string, txId string) (BankCustomerRawData, error)
	TxCreate(ctx context.Context, item BankCustomerRawData, txId string) error
	TxUpdate(ctx context.Context, id string, updateFields UpdateFieldRequests, txId string) error
	TxDelete(ctx context.Context, id string, txId string) error

	// auth
	Authenticate(ctx context.Context, jwt CognitoIDTokenJWT, nowTs int64) BankCustomerAuthResult
}

const TableBankCustomers = "bank_customers"

type BankCustomerRawData struct {
	ID string
	DatabaseItem
}

func (item BankCustomerRawData) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":        item.ID,
		"createdAt": item.CreatedAt,
		"updatedAt": item.UpdatedAt,
	}
}

func newBankCustomerRawDataBySQLRows(rows *sql.Rows) (BankCustomerRawData, error) {
	item := BankCustomerRawData{}
	err := rows.Scan(item.ID, item.CreatedAt, item.UpdatedAt)
	if err != nil {
		return item, crosscuttinglayer.ErrLift(err, context.TODO())
	}
	return item, nil
}

type bankCustomerRepository struct {
	TableRepository[BankCustomerRawData, string]
}

func NewBankCustomerRepository(
	sqlClient SQLClient, txPool TxConnectionPoolIF) BankCustomerRepositoryIF {
	return &bankCustomerRepository{
		TableRepository: TableRepository[BankCustomerRawData, string]{
			tablename:   TableBankCustomers,
			sqlClient:   sqlClient,
			txPool:      txPool,
			convertFunc: newBankCustomerRawDataBySQLRows,
		},
	}
}

func (repo bankCustomerRepository) Authenticate(ctx context.Context, jwt CognitoIDTokenJWT, nowTs int64) BankCustomerAuthResult {
	result := BankCustomerAuthResult{}
	// verify jwt
	header, payload, err := jwt.GetHeaderAndPayload()
	if err != nil {
		result.Err = crosscuttinglayer.ErrLift(err, ctx)
		result.ErrorReason.InternalError = true
		return result
	}
	// kid
	keys, err := repo.getJWKs()
	if err != nil {
		result.Err = crosscuttinglayer.ErrLift(err, ctx)
		result.ErrorReason.InternalError = true
		return result
	}
	if !keys.Have(header.KeyID) {
		appErr := crosscuttinglayer.NewAppUtilError(
			fmt.Sprintf(`header.keyId:%s does not exist in jwks`, header.KeyID),
			ctx,
		)
		result.Err = appErr
		result.ErrorReason.InvalidKid = true
		return result
	}

	// expired
	if !payload.IsValidExp(nowTs) {
		result.Err = crosscuttinglayer.NewAppUtilError("token expired", ctx)
		result.ErrorReason.TokenExpired = true
		return result
	}

	// aud
	if payload.Aud != config.GetCognitoAppClientID() {
		appErr := crosscuttinglayer.NewAppUtilError(
			fmt.Sprintf(`payload.aud:%s is invalid`, payload.Aud),
			ctx,
		)
		result.Err = appErr
		result.ErrorReason.InvalidAud = true
		return result
	}

	// iss
	if !payload.IsValidIss(config.GetCognitoUserPoolID(), config.GetAWSRegion()) {
		appErr := crosscuttinglayer.NewAppUtilError(
			fmt.Sprintf(`payload.iss:%s is invalid`, payload.Iss),
			ctx,
		)
		result.Err = appErr
		result.ErrorReason.InvalidISS = true
		return result
	}

	// token_use
	if !payload.IsIDToken() {
		appErr := crosscuttinglayer.NewAppUtilError(
			fmt.Sprintf(`payload.tokenUse:%s is not id token`, payload.TokenUse),
			ctx,
		)
		result.Err = appErr
		result.ErrorReason.InvalidTokenUse = true
		return result
	}

	customerID := payload.AuthServiceUserName
	result.Customer, err = repo.Get(ctx, customerID)
	if err != nil {
		result.Err = crosscuttinglayer.ErrLift(err, ctx)
		result.ErrorReason.InternalError = true
		return result
	}
	return result
}

func (repo bankCustomerRepository) getJWKs() (keys JSONWebKeys, err error) {
	region := config.GetAWSRegion()
	userPoolID := config.GetCognitoUserPoolID()
	url := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", region, userPoolID)
	response, err := http.Get(url)
	if err != nil {
		err = crosscuttinglayer.NewError(err, context.TODO())
		return keys, err
	}
	defer func() {
		tmpErr := response.Body.Close()
		if tmpErr != nil {
			err = errors.Join(err, tmpErr)
		}
	}()
	var b []byte
	b, err = io.ReadAll(response.Body)
	if err != nil {
		err = crosscuttinglayer.NewError(err, context.TODO())
		return keys, err
	}
	err = json.Unmarshal(b, &keys)
	if err != nil {
		err = crosscuttinglayer.NewError(err, context.TODO())
		return keys, err
	}
	err = nil
	return keys, err
}

type BankCustomerAuthResult struct {
	Customer    BankCustomerRawData
	Err         error
	ErrorReason BankCustomerAuthErrorReason
}

type BankCustomerAuthErrorReason struct {
	TokenExpired    bool
	InvalidKid      bool
	InvalidISS      bool
	InvalidTokenUse bool
	InvalidAud      bool
	InternalError   bool
}
