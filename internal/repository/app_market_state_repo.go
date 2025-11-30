package repository

import (
	"context"

	"victa/internal/domain"
)

type AppMarketStateRepository interface {
	Get(ctx context.Context, appID, marketID int64) (*domain.AppMarketState, error)

	Save(ctx context.Context, s *domain.AppMarketState) (*domain.AppMarketState, error)
}
