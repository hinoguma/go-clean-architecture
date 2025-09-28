package adapter

import (
	"app/applogiclayer/usecase"
	"errors"
)

/********************************
 * Request
 ********************************/
type AddItemToCartRequest struct {
	adapted AdaptedRequest
}

func (r AddItemToCartRequest) UseCaseRequest() usecase.AddItemToCartRequest {
	return usecase.AddItemToCartRequest{
		CustomerUserID: r.adapted.structuredParams["customerUserId"].(string),
		ProductID:      r.adapted.structuredParams["productId"].(string),
		Quantity:       int(r.adapted.structuredParams["quantity"].(float64)),
	}
}

func (r AddItemToCartRequest) Validate() (usecase.AddItemToCartRequest, []error) {
	errs := make([]error, 0)
	ucReq := usecase.AddItemToCartRequest{}
	tmpCuID, ok := r.adapted.GetStructuredParam("customerUserId")
	if ok {
		ucReq.CustomerUserID, ok = tmpCuID.(string)
		if !ok {
			errs = append(errs, errors.New("customerUserId must be string"))
		}
		if ucReq.CustomerUserID == "" {
			errs = append(errs, errors.New("customerUserId is required"))
		}
	} else {
		errs = append(errs, errors.New("customerUserId is required"))
	}

	tmpPdtID, ok := r.adapted.GetStructuredParam("productId")
	if ok {
		ucReq.ProductID, ok = tmpPdtID.(string)
		if !ok {
			errs = append(errs, errors.New("productId must be string"))
		}
		if ucReq.ProductID == "" {
			errs = append(errs, errors.New("productId is required"))
		}
	} else {
		errs = append(errs, errors.New("productId is required"))
	}

	tmpQty, ok := r.adapted.GetStructuredParam("quantity")
	if ok {
		qtyFloat, ok := tmpQty.(float64)
		if !ok {
			errs = append(errs, errors.New("quantity must be number"))
		} else {
			ucReq.Quantity = int(qtyFloat)
			if ucReq.Quantity <= 0 {
				errs = append(errs, errors.New("quantity must be greater than 0"))
			}
		}
	} else {
		errs = append(errs, errors.New("quantity is required"))
	}
	return ucReq, errs
}

/********************************
 * Response
 ********************************/
type AddItemToCartResult struct {
	cartID string
}

/********************************
 * adapter (call usecase)
 ********************************/
type AddItemToCartAdapter interface {
	Execute(reqAdp RequestAdapter, resp AdaptResponder) error
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

func (adp addItemToCartAdapter) Execute(reqAdp RequestAdapter, resp AdaptResponder) error {

	// adapt request
	adaptedReq, err := reqAdp.Do()
	if err != nil {
		return err
	}
	req := AddItemToCartRequest{
		adapted: adaptedReq,
	}

	// validate request
	ucReq, valErrs := req.Validate()
	if len(valErrs) > 0 {
		return errors.Join(valErrs...)
	}

	// call UseCase
	ucResp := adp.uc.Execute(ucReq)

	// adapt response
	adaptedResp := AdapterLayerResponse{
		Body: map[string]string{"cardId": ucResp.CartID},
	}
	resp.Marshal(adaptedResp)

	return nil
}
