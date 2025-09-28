package dip

import (
	"app/crosscutting"
	"app/infrastructure"
)

func NewLogger() crosscutting.Logger {
	return infrastructure.NewCloudWatchLogger()
}
