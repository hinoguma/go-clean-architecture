package apigw

import (
	"app/adapter"
	"app/registory"
	"net/http"
)

func NewServer() *http.ServeMux {

	// prepare adapters
	infraFactory := registory.NewInfraFactory()
	useCaseFactory := registory.NewUseCaseFactory(infraFactory)
	adapterFactory := registory.NewAdapterFactory(
		useCaseFactory,
	)

	sv := http.NewServeMux()
	// GET /account
	sv.Handle("/temp", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle the request for /account
		adp := adapterFactory.NewAddItemToCartAdapter()
		responder := AddItemToCartResponder{w: w}
		err := adp.Execute(RequestAdapter{r: r}, &responder)
		if err != nil {
			return
		}
	}))
	// GET /account/{id}
	return sv
}

type RequestAdapter struct {
	r *http.Request
}

func (adp RequestAdapter) Do() (adapter.AdaptedRequest, error) {
	// e.g. parse body and inject to adapted request
	return adapter.AdaptedRequest{}, nil
}

type AddItemToCartResponder struct {
	w http.ResponseWriter
}

func (responder *AddItemToCartResponder) Marshal(raw adapter.AdapterLayerResponse) {
	responder.w.Write([]byte("ok"))
}
