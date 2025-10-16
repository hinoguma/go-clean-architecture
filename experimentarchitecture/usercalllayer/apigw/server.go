package apigw

import (
	registory2 "app/experimentarchitecture/registory"
	"app/experimentarchitecture/usercalllayer/apigw/handler"
	"net/http"
)

func NewServer() *http.ServeMux {

	// prepare adapters
	infraFactory := registory2.NewInfraFactory()
	useCaseFactory := registory2.NewUseCaseFactory(infraFactory)
	adapterFactory := registory2.NewAdapterFactory(
		useCaseFactory,
	)

	sv := http.NewServeMux()
	// GET /account
	sv.Handle("/temp", handler.AddItemToCartHandler(adapterFactory))
	// GET /account/{id}
	return sv
}
