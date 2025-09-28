package model

type ShoppingCartID string

type CartItem struct {
	ProductID ProductItemInfoID `json:"product_id"`
	Quantity  int               `json:"quantity"`
}

type ShoppingCart struct {
	ID             ShoppingCartID `json:"id"`
	CustomerUserID CustomerUserID `json:"customer_user_id"`
	Items          []CartItem     `json:"items"`
	CreatedAt      float64        `json:"createdAt"`
	UpdatedAt      float64        `json:"updatedAt"`
}
