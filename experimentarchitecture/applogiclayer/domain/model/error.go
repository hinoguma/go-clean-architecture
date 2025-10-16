package model

import (
	"app/experimentarchitecture/crosscutting/utils"
)

type ValidateError struct {
	utils.ValidateErrorDetail
}
