package dip

import (
	"app/crosscutting"
	"app/infrastructure/log"
)

func NewLogger() crosscutting.Logger {
	return log.NewCloudWatchLogger()
}
