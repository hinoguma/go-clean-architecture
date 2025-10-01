package apigw

import (
	"app/registory"
	"app/userinterfacelayer/apigw/handler"
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
	sv.Handle("/temp", handler.AddItemToCartHandler(adapterFactory))
	// GET /account/{id}
	return sv
}
