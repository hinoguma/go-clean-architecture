package model




import "time"


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