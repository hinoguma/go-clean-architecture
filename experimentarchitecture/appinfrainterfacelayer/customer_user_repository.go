package appinfrainterfacelayer

import (
	"app/experimentarchitecture/applogiclayer/domain/model"
	"context"
)

type CustomerUserRepository interface {
	Get(ctx context.Context, id model.CustomerUserID) (model.CustomerUser, error)
}
