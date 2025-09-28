package additemtocart

import (
	"app/adapter"
	"app/registory"
	"flag"
	"fmt"
)

func main() {
	// e.g. parse args from CLI
	args := []string{"--user", "user1", "--item", "itemA", "--quantity", "2"}
	// prepare adapters
	infraFactory := registory.NewInfraFactory()
	useCaseFactory := registory.NewUseCaseFactory(infraFactory)
	adapterFactory := registory.NewAdapterFactory(
		useCaseFactory,
	)

	adp := adapterFactory.NewAddItemToCartAdapter()
	responder := AddItemToCartResponder{}
	err := adp.Execute(RequestAdapter{args: args}, &responder)
	if err != nil {
		return
	}
}

type RequestAdapter struct {
	args []string
}

func (adp RequestAdapter) Do() (adapter.AdaptedRequest, error) {
	// e.g. parse body and inject to adapted request
	cuidPtr := flag.String("cuid", "", "userId")
	pidPtr := flag.String("pid", "", "productId")
	qPtr := flag.Int("quantity", 0, "quantity")
	flag.Parse()

	req := adapter.AdaptedRequest{}
	req.AddStructuredParam("campaignUserId", *cuidPtr)
	req.AddStructuredParam("productId", *pidPtr)
	req.AddStructuredParam("quantity", *qPtr)
	return req, nil

}

type AddItemToCartResponder struct {
}

func (responder *AddItemToCartResponder) Marshal(raw adapter.AdapterLayerResponse) {
	fmt.Sprintln("result is", raw.Body.(string))
}
