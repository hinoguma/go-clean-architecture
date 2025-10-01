package adapter

type AdaptedRequest struct {
	From               RequestFrom
	requestID          string
	requestIP          string
	authorizationToken string
	unstructuredParams string
	structuredParams   map[string]interface{}
	files              map[string][]byte

	// HTTP Request
	HTTPRequest HTTPRequest

	// Command Line Interface
	CLIArguments []CLIArgument
}

func (req AdaptedRequest) IsHttpRequest() bool {
	return req.IsHTTPServerRequest() || req.IsAPIGatewayRequest()
}

func (req AdaptedRequest) IsAPIGatewayRequest() bool {
	return req.From == RequestFromAPIGateway
}

func (req AdaptedRequest) IsHTTPServerRequest() bool {
	return req.From == RequestFromHTTPServer
}

func (req AdaptedRequest) IsCLIRequest() bool {
	return req.From == RequestFromCLI
}

func (req AdaptedRequest) GetStructuredParam(key string) (interface{}, bool) {
	v, ok := req.structuredParams[key]
	return v, ok
}

func (req AdaptedRequest) AddStructuredParam(key string, value interface{}) {
	if req.structuredParams == nil {
		req.structuredParams = make(map[string]interface{})
	}
	req.structuredParams[key] = value
}

type RequestAdapter interface {
	Do() (AdaptedRequest, error)
}

//type APIGatewayProxyRequest struct {
//	Resource                        string                        `json:"resource"` // The resource path defined in API Gateway
//	Path                            string                        `json:"path"`     // The url path for the caller
//	HTTPMethod                      string                        `json:"httpMethod"`
//	Headers                         map[string]string             `json:"headers"`
//	MultiValueHeaders               map[string][]string           `json:"multiValueHeaders"`
//	QueryStringParameters           map[string]string             `json:"queryStringParameters"`
//	MultiValueQueryStringParameters map[string][]string           `json:"multiValueQueryStringParameters"`
//	PathParameters                  map[string]string             `json:"pathParameters"`
//	StageVariables                  map[string]string             `json:"stageVariables"`
//	RequestContext                  APIGatewayProxyRequestContext `json:"requestContext"`
//	Body                            string                        `json:"body"`
//	IsBase64Encoded                 bool                          `json:"isBase64Encoded,omitempty"`
//}
//type APIGatewayProxyRequestContext struct {
//	AccountID         string                    `json:"accountId"`
//	ResourceID        string                    `json:"resourceId"`
//	OperationName     string                    `json:"operationName,omitempty"`
//	Stage             string                    `json:"stage"`
//	DomainName        string                    `json:"domainName"`
//	DomainPrefix      string                    `json:"domainPrefix"`
//	RequestID         string                    `json:"requestId"`
//	ExtendedRequestID string                    `json:"extendedRequestId"`
//	Protocol          string                    `json:"protocol"`
//	Identity          APIGatewayRequestIdentity `json:"identity"`
//	ResourcePath      string                    `json:"resourcePath"`
//	Path              string                    `json:"path"`
//	Authorizer        map[string]infrainterface{}    `json:"authorizer"`
//	HTTPMethod        string                    `json:"httpMethod"`
//	RequestTime       string                    `json:"requestTime"`
//	RequestTimeEpoch  int64                     `json:"requestTimeEpoch"`
//	APIID             string                    `json:"apiId"` // The API Gateway rest API Id
//}
//type APIGatewayRequestIdentity struct {
//	CognitoIdentityPoolID         string `json:"cognitoIdentityPoolId"`
//	AccountID                     string `json:"accountId"`
//	CognitoIdentityID             string `json:"cognitoIdentityId"`
//	Caller                        string `json:"caller"`
//	APIKey                        string `json:"apiKey"`
//	APIKeyID                      string `json:"apiKeyId"`
//	AccessKey                     string `json:"accessKey"`
//	SourceIP                      string `json:"sourceIp"`
//	CognitoAuthenticationType     string `json:"cognitoAuthenticationType"`
//	CognitoAuthenticationProvider string `json:"cognitoAuthenticationProvider"`
//	UserArn                       string `json:"userArn"` //nolint: stylecheck
//	UserAgent                     string `json:"userAgent"`
//	User                          string `json:"user"`
//}
//
//type APIGWProxyRequestAdapter struct {
//	ctx context.Context
//	raw APIGatewayProxyRequest
//}
//
//func (adapter APIGWProxyRequestAdapter) do() (AdaptedRequest, error)  {
//	ad := AdaptedRequest{}
//	sp := make(map[string]infrainterface{})
//	for k, v := range adapter.raw.PathParameters {
//		sp[k] = v
//	}
//	for k, v := range adapter.raw.QueryStringParameters {
//		sp[k] = v
//	}
//	for k, v := range adapter.raw.MultiValueQueryStringParameters {
//		sp[k] = v
//	}
//	b := []byte(adapter.raw.Body)
//	jm := make(map[string]infrainterface{})
//	err := json.Unmarshal(b, &jm)
//	if err != nil {
//		return ad, err
//	}
//	for k, v := range jm {
//		sp[k] = v
//	}
//	ad.structuredParams = sp
//	return ad, nil
//}
//
//
//type CommandLineInterfaceAdapter struct {
//	ctx context.Context
//}
//
//func (adapter CommandLineInterfaceAdapter) do() (AdaptedRequest, error)  {
//	ad := AdaptedRequest{}
//	sp := make(map[string]infrainterface{})
//
//
//	return ad, nil
//}

type RequestFrom string

const (
	RequestFromAPIGateway RequestFrom = "APIGateway"
	RequestFromHTTPServer RequestFrom = "HTTPServer"
	RequestFromCLI        RequestFrom = "CLI"
)

type HTTPRequest struct {
	Method  string
	Headers map[string]string
	Body    string

	QueryParams map[string]string
	PathParams  map[string]string
}

type CLIArgument struct {
	Name  string
	Value string
}
