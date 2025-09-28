package infrainterface

import (
	"app/applogiclayer/domain/model"
	"context"
)

type ShoppingCartRepository interface {
	GetByUserID(ctx context.Context, userID model.CustomerUserID) (model.ShoppingCart, error)
}
