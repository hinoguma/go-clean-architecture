package handler

import (
	"app/adapter"
	"net/http"
)

type RequestAdapter struct {
	r *http.Request
}

func (adp RequestAdapter) Do() (adapter.AdaptedRequest, error) {
	// e.g. parse body and inject to adapted request
	return adapter.AdaptedRequest{}, nil
}
