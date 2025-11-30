package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"victa/internal/domain"
)

type AppMarketStateRepo struct {
	db       *sql.DB
	stGet    *sql.Stmt
	stUpsert *sql.Stmt
}

func NewAppMarketStateRepo(db *sql.DB) (*AppMarketStateRepo, error) {
	r := &AppMarketStateRepo{db: db}
	var err error

	if r.stGet, err = db.Prepare(`
		SELECT app_id,
		       market_id,
		       last_version,
		       last_review_id,
		       updated_at
		  FROM app_market_states
		 WHERE app_id = $1 AND market_id = $2
	`); err != nil {
		return nil, fmt.Errorf("prepare get app_market_state: %w", err)
	}

	if r.stUpsert, err = db.Prepare(`
		INSERT INTO app_market_states (app_id, market_id, last_version, last_review_id, updated_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (app_id, market_id)
		DO UPDATE SET
			last_version   = EXCLUDED.last_version,
			last_review_id = EXCLUDED.last_review_id,
			updated_at     = NOW()
		RETURNING app_id, market_id, last_version, last_review_id, updated_at
	`); err != nil {
		return nil, fmt.Errorf("prepare upsert app_market_state: %w", err)
	}

	return r, nil
}

func (r *AppMarketStateRepo) Close() error {
	for _, st := range []*sql.Stmt{r.stGet, r.stUpsert} {
		if st != nil {
			if err := st.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *AppMarketStateRepo) Get(
	ctx context.Context,
	appID, marketID int64,
) (*domain.AppMarketState, error) {

	var (
		lastVersion   sql.NullInt64
		lastReviewStr sql.NullString
		updatedAt     time.Time
	)

	err := r.stGet.QueryRowContext(ctx, appID, marketID).
		Scan(&appID, &marketID, &lastVersion, &lastReviewStr, &updatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get app_market_state: %w", err)
	}

	state := &domain.AppMarketState{
		AppID:     appID,
		MarketID:  marketID,
		UpdatedAt: updatedAt,
	}

	if lastVersion.Valid {
		v := lastVersion.Int64
		state.LastVersion = &v
	}
	if lastReviewStr.Valid {
		s := lastReviewStr.String
		state.LastReviewID = &s
	}

	return state, nil
}

func (r *AppMarketStateRepo) Save(
	ctx context.Context,
	s *domain.AppMarketState,
) (*domain.AppMarketState, error) {

	var (
		lastVersion   sql.NullInt64
		lastReviewStr sql.NullString
		updatedAt     time.Time
	)

	if s.LastVersion != nil {
		lastVersion = sql.NullInt64{Int64: *s.LastVersion, Valid: true}
	}
	if s.LastReviewID != nil {
		lastReviewStr = sql.NullString{String: *s.LastReviewID, Valid: true}
	}

	err := r.stUpsert.QueryRowContext(
		ctx,
		s.AppID,
		s.MarketID,
		lastVersion,
		lastReviewStr,
	).Scan(
		&s.AppID,
		&s.MarketID,
		&lastVersion,
		&lastReviewStr,
		&updatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("upsert app_market_state: %w", err)
	}

	s.UpdatedAt = updatedAt

	if lastVersion.Valid {
		v := lastVersion.Int64
		s.LastVersion = &v
	} else {
		s.LastVersion = nil
	}
	if lastReviewStr.Valid {
		str := lastReviewStr.String
		s.LastReviewID = &str
	} else {
		s.LastReviewID = nil
	}

	return s, nil
}
