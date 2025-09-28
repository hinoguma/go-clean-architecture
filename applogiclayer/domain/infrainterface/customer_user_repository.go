package infrainterface

import (
	"app/applogiclayer/domain/model"
	"context"
)

type CustomerUserRepository interface {
	Get(ctx context.Context, id model.CustomerUserID) (model.CustomerUser, error)
}
