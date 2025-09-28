package infrainterface

import (
	"app/applogiclayer/domain/model"
	"context"
)

type ProductItemInfoRepository interface {
	Get(ctx context.Context, id model.ProductItemInfoID) (model.ProductItemInfo, error)
}
