package model

import "app/threelayeredarchitecture/appinfraadapterlayer"

type BankCustomer struct {
	ID   string
	Name string

	DataItem
}

func (model BankCustomer) SetFromDTO(dto appinfraadapterlayer.BankCustomerDTO) {
	model.ID = dto.ID
	model.Name = dto.Name
	model.DataItem.SetFromDTO(dto.DatabaseItem)
}
