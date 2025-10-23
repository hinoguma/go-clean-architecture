package appinfraadapterlayer

import "context"

type AuthenticateRequestDTO struct {
	Token string
}

type AuthenticateResultDTO struct {
	Success     bool
	Err         error
	ErrorReason AuthErrorReason
}

type AuthClientAdapterIF interface {
	Authenticate(ctx context.Context, req AuthenticateRequestDTO) AuthenticateResultDTO
}
