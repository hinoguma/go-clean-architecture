package infra

import (
	"app/threelayeredarchitecture/crosscuttinglayer/adapter"
	"github.com/google/uuid"
)

type UUIDV4Generator struct {
}

func (generator UUIDV4Generator) Issue() string {
	return uuid.New().String()
}

func NewUUIDV4Generator() adapter.StrIDGenerator {
	return UUIDV4Generator{}
}
