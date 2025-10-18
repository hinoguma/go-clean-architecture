package model

import (
	"app/threelayeredarchitecture/appinfraadapterlayer"
	"time"
)


type HasCreatedAt struct {
	CreatedAt time.Time
}

type HasUpdatedAt struct {
	UpdatedAt time.Time
}

type DataItem struct {
	HasCreatedAt
	HasUpdatedAt
}

func (model *DataItem) SetFromDTO(dto appinfraadapterlayer.DatabaseItem) {
	model.CreatedAt = time.Unix(dto.CreatedAt, 0)
	model.UpdatedAt = time.Unix(dto.UpdatedAt, 0)
}