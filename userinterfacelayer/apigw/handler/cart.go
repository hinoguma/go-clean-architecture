package handler

import (
	"app/crosscutting/logger"
	"app/crosscutting/utils"
	"app/registory"
	"net/http"
)

func AddItemToCartHandler(
	adapterFactory registory.AdapterFactory,
) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// Handle the request for /account
			ctx := r.Context()
			adp := adapterFactory.NewAddItemToCartAdapter()
			req, err := RequestAdapter{r: r}.Do()
			if err != nil {
				logger.LogError(r.Context(), utils.NewLogRequest(err.Error()).ToValue())
				w.Write([]byte("error"))
				return
			}
			resp, err := adp.Execute(ctx, req)
			if err != nil {
				logger.LogError(r.Context(), utils.NewLogRequest(err.Error()).ToValue())
				w.Write([]byte("error"))
				return
			}
			w.Write([]byte("ok: " + resp.CartID))
		},
	)
}
