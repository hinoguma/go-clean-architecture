package model

type CustomerUserID string

type CustomerUser struct {
	ID   CustomerUserID `json:"id"`
	Name string         `json:"name"`

	CreatedAt float64 `json:"createdAt"`
	UpdatedAt float64 `json:"updatedAt"`
}
