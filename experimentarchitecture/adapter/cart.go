package adapter

import (
	"app/experimentarchitecture/applogiclayer/usecase"
	"context"
	"fmt"
)

/********************************
 * Response
 ********************************/
type AddItemToCartResult struct {
	CartID string
	BaseResult
}

func NewAddItemToCartResult(ucRes usecase.AddItemToCartResult) AddItemToCartResult {
	res := AddItemToCartResult{
		CartID: ucRes.CartID,
	}
	res.SetByBaseUseCaseResult(ucRes.BaseUseCaseResult)
	return res
}

/********************************
 * adapter (call usecase)
 ********************************/
type AddItemToCartAdapter interface {
	Execute(ctx context.Context, req AdaptedRequest) (AddItemToCartResult, error)
}

func NewAddItemToCartAdapter(
	uc usecase.AddItemToCartUseCase,
) AddItemToCartAdapter {
	return addItemToCartAdapter{
		uc: uc,
	}
}

type addItemToCartAdapter struct {
	// auth service
	uc usecase.AddItemToCartUseCase
}

func (adp addItemToCartAdapter) validateRequest(req AdaptedRequest) (usecase.AddItemToCartRequest, []ValidateError) {
	vErrs := make([]ValidateError, 0)
	ucReq := usecase.AddItemToCartRequest{}

	if req.IsHttpRequest() {
		return adp.validateRequestForHTTP(req)
	}
	if req.IsCLIRequest() {
		return adp.validateRequestForCLI(req)
	}
	// unknown request type
	vErrs = append(vErrs, NewValidateError("", fmt.Sprintf("unknown request type. requestFrom:%s", req.From)))
	return ucReq, vErrs
}

func (adp addItemToCartAdapter) validateRequestForHTTP(req AdaptedRequest) (usecase.AddItemToCartRequest, []ValidateError) {
	vErrs := make([]ValidateError, 0)
	ucReq := usecase.AddItemToCartRequest{}
	tmpCuID, ok := req.GetStructuredParam("customerUserId")
	if ok {
		ucReq.CustomerUserID, ok = tmpCuID.(string)
		if !ok {
			vErrs = append(vErrs, NewValidateErrorRequired("customerUserId"))
		}
		if ucReq.CustomerUserID == "" {
			vErrs = append(vErrs, NewValidateErrorRequired("customerUserId"))
		}
	} else {
		vErrs = append(vErrs, NewValidateErrorRequired("customerUserId"))
	}

	tmpPdtID, ok := req.GetStructuredParam("productId")
	if ok {
		ucReq.ProductID, ok = tmpPdtID.(string)
		if !ok {
			vErrs = append(vErrs, NewValidateErrorRequired("productId"))
		} else {
			if ucReq.ProductID == "" {
				vErrs = append(vErrs, NewValidateErrorRequired("productId"))
			}
		}
	} else {
		vErrs = append(vErrs, NewValidateErrorRequired("productId"))
	}

	tmpQty, ok := req.GetStructuredParam("quantity")

	if ok {
		qty, ok := tmpQty.(float64)
		if !ok {
			vErrs = append(vErrs, NewValidateErrorRequired("quantity"))
		} else {
			ucReq.Quantity = int(qty)
			if ucReq.Quantity <= 0 {
				vErrs = append(vErrs, NewValidateError("quantity", "must be greater than zero"))
			}
		}
	} else {
		vErrs = append(vErrs, NewValidateErrorRequired("quantity"))
	}
	return ucReq, vErrs
}

func (adp addItemToCartAdapter) validateRequestForCLI(req AdaptedRequest) (usecase.AddItemToCartRequest, []ValidateError) {
	vErrs := make([]ValidateError, 0)
	vErrs = append(vErrs, NewValidateError("", "CLI not supported"))
	return usecase.AddItemToCartRequest{}, vErrs
}

func (adp addItemToCartAdapter) Execute(ctx context.Context, req AdaptedRequest) (AddItemToCartResult, error) {

	// validate request
	ucReq, valErrs := adp.validateRequest(req)
	if len(valErrs) > 0 {
		return AddItemToCartResult{}, NewAdapterErrorValidation(valErrs)
	}

	// call UseCase
	ucResp, ucErr := adp.uc.Execute(ucReq)
	if ucErr != nil {
		return AddItemToCartResult{}, handleErrFromUseCase(ucErr)
	}
	res := NewAddItemToCartResult(ucResp)
	res.SetByCtx(ctx)
	if !res.IsSuccess() {
		return res, NewResultNotSuccessError(res.BaseResult)
	}
	return res, nil
}
