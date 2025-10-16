package appinfrainterfacelayer

import (
	model2 "app/experimentarchitecture/applogiclayer/domain/model"
	"context"
)

type ShoppingCartRepository interface {
	GetByUserID(ctx context.Context, userID model2.CustomerUserID) (model2.ShoppingCart, error)
}
