package transfer

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "net/http"

    "app/threelayerdarchitecture/applogic/domain"
    "app/threelayerdarchitecture/applogic/usecase"
)

type HTTPRequest struct {
    Path          string
    Body          []byte
    Headers       map[string]string
    RequestUserID string
    CorrelationID string
}

type HTTPResponse struct {
    StatusCode int
    Body       string
    Headers    map[string]string
}

var (
    ErrMissingJSONBody   = errors.New("transfer adapter: body is required")
    ErrMissingAuthHeader = errors.New("transfer adapter: user header missing")
)

type Adapter interface {
    Execute(ctx context.Context, req HTTPRequest) (HTTPResponse, error)
}

type transferAdapter struct {
    usecase usecase.TransferUseCase
}

func NewTransferAdapter(uc usecase.TransferUseCase) Adapter {
    return transferAdapter{usecase: uc}
}

func (a transferAdapter) Execute(ctx context.Context, req HTTPRequest) (HTTPResponse, error) {
    if len(req.Body) == 0 {
        return errorResponse(http.StatusBadRequest, "missing_body"), ErrMissingJSONBody
    }

    userID := req.RequestUserID
    if userID == "" {
        if val, ok := req.Headers["X-User-ID"]; ok {
            userID = val
        } else if val, ok := req.Headers["x-user-id"]; ok {
            userID = val
        }
    }
    if userID == "" {
        return errorResponse(http.StatusUnauthorized, "missing_user"), ErrMissingAuthHeader
    }

    var payload struct {
        FromAccountID string `json:"from_account_id"`
        ToAccountID   string `json:"to_account_id"`
        Amount        int64  `json:"amount"`
    }
    if err := json.Unmarshal(req.Body, &payload); err != nil {
        return errorResponse(http.StatusBadRequest, "invalid_json"), fmt.Errorf("decode request: %w", err)
    }

    useReq := usecase.TransferRequest{
        InitiatingUserID: userID,
        FromAccountID:    payload.FromAccountID,
        ToAccountID:      payload.ToAccountID,
        Amount:           payload.Amount,
        CorrelationID:    req.CorrelationID,
    }

    result, err := a.usecase.Execute(ctx, useReq)
    if err != nil {
        return mapError(err), err
    }

    responseBody, _ := json.Marshal(map[string]string{
        "transaction_id": result.TransactionID,
    })

    return HTTPResponse{
        StatusCode: http.StatusOK,
        Body:       string(responseBody),
        Headers: map[string]string{
            "Content-Type": "application/json",
        },
    }, nil
}

func errorResponse(code int, message string) HTTPResponse {
    body, _ := json.Marshal(map[string]string{"error": message})
    return HTTPResponse{StatusCode: code, Body: string(body), Headers: map[string]string{"Content-Type": "application/json"}}
}

func mapError(err error) HTTPResponse {
    switch {
    case errors.Is(err, usecase.ErrInvalidTransferRequest):
        return errorResponse(http.StatusBadRequest, "invalid_request")
    case errors.Is(err, usecase.ErrUnauthorizedUser):
        return errorResponse(http.StatusForbidden, "unauthorized")
    case errors.Is(err, domain.ErrAmountNegative):
        return errorResponse(http.StatusBadRequest, "invalid_amount")
    case errors.Is(err, domain.ErrInsufficientBalance):
        return errorResponse(http.StatusConflict, "insufficient_funds")
    default:
        return errorResponse(http.StatusInternalServerError, "transfer_failed")
    }
}
