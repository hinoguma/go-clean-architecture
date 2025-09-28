package adapter

type AdaptResponder interface {
	Marshal(raw AdapterLayerResponse)
}

type AdapterLayerResponse struct {
	Type  ResultType
	Body  interface{}
	Error Error
}

func (res AdapterLayerResponse) IsSuccess() bool {
	return res.Type == ResultTypeSuccess || res.Type == ResultTypeSuccessAsync
}

func (res AdapterLayerResponse) IsFailure() bool {
	return !res.IsSuccess()
}

type ResultType string

const (
	ResultTypeSuccess      ResultType = "success"
	ResultTypeSuccessAsync ResultType = "success_async"
	ResultTypeError        ResultType = "error"
)
