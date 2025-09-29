package crosscutting

import (
	"app/crosscutting/infrainterface"
	"context"
	"encoding/json"
	"fmt"
)

func SetGlobalLogger(logger infrainterface.Logger) {
	globalLogger = logger
}

var globalLogger infrainterface.Logger = NewStdLogger()

func LogInfo(ctx context.Context, req infrainterface.LogRequest) {
	if globalLogger == nil {
		return
	}
	globalLogger.Info(ctx, req)
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

func (req LogRequest) LogJson(level string) infrainterface.LogJson {
	return LogJson{
		RequestID:  req.RequestID,
		Message:    req.Message,
		Level:      level,
		Time:       req.Time,
		Additional: req.Additional,
	}
}

type LogJson struct {
	RequestID  string                 `json:"requestId,omitempty"`
	Message    string                 `json:"message"`
	Level      string                 `json:"level"`
	Time       string                 `json:"time"`
	Additional map[string]interface{} `json:"additional,omitempty"`
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

type StdLogger struct {
}

func NewStdLogger() infrainterface.Logger {
	return StdLogger{}
}

func (s StdLogger) Info(ctx context.Context, req infrainterface.LogRequest) {
	logjson := req.LogJson("INFO")
	fmt.Println(logjson.JsonString())
}
func (s StdLogger) Error(ctx context.Context, req infrainterface.LogRequest) {
	logjson := req.LogJson("ERROR")
	fmt.Println(logjson.JsonString())
}

func (s StdLogger) Debug(ctx context.Context, req infrainterface.LogRequest) {
	logjson := req.LogJson("DEBUG")
	fmt.Println(logjson.JsonString())
}
