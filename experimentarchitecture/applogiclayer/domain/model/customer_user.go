package model

import (
	"app/experimentarchitecture/crosscutting/utils"
)

type CustomerUserID string

func (value CustomerUserID) String() string {
	return string(value)
}

type CustomerUser struct {
	ID   CustomerUserID `json:"id"`
	Name string         `json:"name"`

	CreatedAt utils.UnixTimestamp `json:"createdAt"`
	UpdatedAt utils.UnixTimestamp `json:"updatedAt"`
}
