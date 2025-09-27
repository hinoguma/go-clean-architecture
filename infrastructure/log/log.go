package log

import "app/crosscutting"

type CloudWatchLogger struct {
}

func (logger CloudWatchLogger) Info(req crosscutting.LogRequest) {
	//TODO implement me
	panic("implement me")
}

func (logger CloudWatchLogger) Error(req crosscutting.LogRequest) {
	//TODO implement me
	panic("implement me")
}

func (logger CloudWatchLogger) Debug(req crosscutting.LogRequest) {
	//TODO implement me
	panic("implement me")
}

func NewCloudWatchLogger() crosscutting.Logger {
	return CloudWatchLogger{}
}
