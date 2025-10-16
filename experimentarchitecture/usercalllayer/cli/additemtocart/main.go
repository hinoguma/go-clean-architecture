package additemtocart

import (
	"app/experimentarchitecture/adapter"
	"app/experimentarchitecture/crosscutting/logger"
	"app/experimentarchitecture/crosscutting/utils"
	registory2 "app/experimentarchitecture/registory"
	"context"
	"flag"
	"fmt"
)

func main() {
	// e.g. parse args from CLI
	ctx := context.Background()

	// prepare adapters
	infraFactory := registory2.NewInfraFactory()
	useCaseFactory := registory2.NewUseCaseFactory(infraFactory)
	adapterFactory := registory2.NewAdapterFactory(
		useCaseFactory,
	)

	adp := adapterFactory.NewAddItemToCartAdapter()
	adaptedReq, err := AdaptedRequest()
	if err != nil {
		return
	}
	resp, err := adp.Execute(ctx, adaptedReq)
	if err != nil {
		logger.LogError(ctx, utils.NewLogRequest(err.Error()).ToValue())
		return
	}
	fmt.Sprintf("resp: %+v", resp)
}

func AdaptedRequest() (adapter.AdaptedRequest, error) {
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
