package crosscutting

import (
	"app/crosscutting/dip"
)

var globalLogger = dip.NewLogger()

func LogInfo(req LogRequest) {
	globalLogger.Info(req)
}

// implements in infrastructure layer
type Logger interface {
	Info(req LogRequest)
	Error(req LogRequest)
	Debug(req LogRequest)
}

type LogRequest struct {
	RequestID  string                 `json:"requestId,omitempty"`
	Message    string                 `json:"message"`
	Time       string                 `json:"time"`
	Additional map[string]interface{} `json:"additional,omitempty"`
}

func NewLogRequest(message string) *LogRequest {
	return &LogRequest{
		RequestID:  "",
		Message:    message,
		Time:       "",
		Additional: nil,
	}
}

func (req *LogRequest) WithRequestID(requestID string) *LogRequest {
	req.RequestID = requestID
	return req
}

func (req *LogRequest) WithMessage(message string) *LogRequest {
	req.Message = message
	return req
}

func (req *LogRequest) WithTime(time string) *LogRequest {
	req.Time = time
	return req
}

func (req *LogRequest) AdditionalInfo(key string, value interface{}) *LogRequest {
	if req.Additional == nil {
		req.Additional = make(map[string]interface{})
	}
	req.Additional[key] = value
	return req
}

func (req *LogRequest) ToValue() LogRequest {
	return *req
}

type LogJson struct {
	RequestID  string                 `json:"requestId,omitempty"`
	Message    string                 `json:"message"`
	Level      string                 `json:"level"`
	Time       string                 `json:"time"`
	Additional map[string]interface{} `json:"additional,omitempty"`
}
