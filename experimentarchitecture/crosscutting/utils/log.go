package utils

import (
	"context"
	"encoding/json"
	"fmt"
)

type LogLevel string

func (v LogLevel) String() string {
	return string(v)
}

const (
	LogLevelInfo  LogLevel = "INFO"
	LogLevelError LogLevel = "ERROR"
	LogLevelDebug LogLevel = "DEBUG"
	LogLevelFatal LogLevel = "FATAL"
)

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

func (req LogRequest) LogJson() LogJson {
	return LogJson{
		RequestID:  req.RequestID,
		Message:    req.Message,
		Time:       req.Time,
		Additional: req.Additional,
	}
}

type LogJson struct {
	RequestID  string                 `json:"requestId,omitempty"`
	Message    string                 `json:"message"`
	Level      LogLevel               `json:"level"`
	Time       string                 `json:"time"`
	Additional map[string]interface{} `json:"additional,omitempty"`
}

func (lj *LogJson) WithRequestIDByCtx(ctx context.Context) *LogJson {
	lj.RequestID = GetRequestID(ctx)
	return lj
}

func (lj *LogJson) WithLevel(level LogLevel) *LogJson {
	lj.Level = level
	return lj
}

func (lj LogJson) JsonString() string {
	b, err := json.Marshal(lj)
	if err != nil {
		return fmt.Sprintf(
			`{"requestId":"%s","message":"%s","level":"%s","time":"%s","additional":%v,"logError":"%s"}`,
			lj.RequestID,
			lj.Message,
			lj.Level,
			lj.Time,
			lj.Additional,
			err.Error(),
		)
	}
	return string(b)
}
