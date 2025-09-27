package model

import (
	"app/crosscutting/utils"
)

type DataItemCommonProps struct {
	CreatedAt utils.UnixTimestamp `json:"createdAt"`
	UpdatedAt utils.UnixTimestamp `json:"updatedAt"`
}
