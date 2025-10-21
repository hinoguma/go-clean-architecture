package appinfraadapterlayer

type DatabaseItem struct {
	CreatedAt int64
	UpdatedAt int64
}

type UpdateDatabaseItemRequestDTO struct {
	UpdatedAt *int64
}

func (dto *UpdateDatabaseItemRequestDTO) WithUpdatedAt(val int64) *UpdateDatabaseItemRequestDTO {
	dto.UpdatedAt = &val
	return dto
}
