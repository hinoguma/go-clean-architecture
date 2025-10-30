package infra

import (
	"github.com/google/uuid"
)

type UUIDV4Generator struct {
}

func (generator UUIDV4Generator) Issue() string {
	return uuid.New().String()
}
