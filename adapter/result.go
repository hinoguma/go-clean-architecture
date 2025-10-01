package adapter

import (
	"app/applogiclayer/usecase"
	"app/crosscutting/utils"
	"context"
	"fmt"
)

type BaseResult struct {
	RequestID  string
	ResultType ResultType
}

type ResultType string

const (
	ResultTypeSuccessCompleted ResultType = "SUCCESSFULLY_COMPLETED"
	ResultTypeSuccessAccepted  ResultType = "SUCCESSFULLY_ACCEPTED"
)

func (res *BaseResult) SetByBaseUseCaseResult(ucRes usecase.BaseUseCaseResult) *BaseResult {
	if ucRes.IsSuccessCompleted() {
		res.ResultType = ResultTypeSuccessCompleted
	}
	if ucRes.IsSuccessAccepted() {
		res.ResultType = ResultTypeSuccessAccepted
	}
	return res
}

func (res *BaseResult) SetByCtx(ctx context.Context) *BaseResult {
	if ctx == nil {
		return res
	}
	requestID := utils.GetRequestID(ctx)
	res.RequestID = requestID
	return res
}

func (res BaseResult) IsSuccess() bool {
	return res.IsSuccessCompleted() || res.IsSuccessAccepted()
}

func (res BaseResult) IsSuccessCompleted() bool {
	return res.ResultType == ResultTypeSuccessCompleted
}

func (res BaseResult) IsSuccessAccepted() bool {
	return res.ResultType == ResultTypeSuccessAccepted
}

func NewResultNotSuccessError(res BaseResult) AdapterError {
	return AdapterError{
		ErrType: ErrorTypeApplication,
		Err: utils.NewError(
			fmt.Sprintf("ResultType is not success. ResultType: %s", res.ResultType),
		),
		ValidateErrors: nil,
	}
}
