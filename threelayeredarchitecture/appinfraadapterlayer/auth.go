package appinfraadapterlayer

import "context"

type AuthenticateRequestDTO struct {
	Token     string
	Timestamp int64
}

type AuthenticateResultDTO struct {
	Success     bool
	Err         error
	ErrorReason AuthErrorReason
}

type AuthClientAdapterIF interface {
	Authenticate(ctx context.Context, req AuthenticateRequestDTO) AuthenticateResultDTO
}
