package utils

type APIStatusCode int

// ref: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status
const (
	// 1xx: Informational

	// 2xx: Success
	APIStatusOK        APIStatusCode = 200
	APIStatusAccepted  APIStatusCode = 202
	APIStatusNoContent APIStatusCode = 204

	// 3xx: Redirection
	APIStatusNotModified APIStatusCode = 304

	// 4xx: Client Errors
	APIStatusBadRequest      APIStatusCode = 400
	APIStatusUnauthorized    APIStatusCode = 401
	APIStatusForbidden       APIStatusCode = 403
	APIStatusRequestTimeout  APIStatusCode = 408
	APIStatusPayloadTooLarge APIStatusCode = 413
	APIStatusTooManyRequests APIStatusCode = 429
	APIStatusNotFound        APIStatusCode = 404
	APIStatusConflict        APIStatusCode = 409

	// 5xx: Server Errors
	APIStatusInternalServerError APIStatusCode = 500
	APIStatusBadGateway          APIStatusCode = 502
	APIStatusServiceUnavailable  APIStatusCode = 503
	APIStatusGatewayTimeout      APIStatusCode = 504
)

type APIErrorResponse struct {
	Error APIError
}

type APIErrorType string

type APIError struct {
	Type    APIErrorType
	Message string
	Details []ValidateErrorDetail
}
