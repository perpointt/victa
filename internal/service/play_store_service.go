package service

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"time"

	"victa/internal/domain"
	"victa/internal/repository"

	"google.golang.org/api/androidpublisher/v3"
	"google.golang.org/api/option"
)

var playNameRe = regexp.MustCompile(`^(\d+)\s+\((\d+\.\d+\.\d+)\)$`)

type PlayStoreService struct {
	pub       *androidpublisher.Service
	edits     *androidpublisher.EditsService
	reviews   *androidpublisher.ReviewsService
	tracks    *androidpublisher.EditsTracksService
	stateRepo repository.AppMarketStateRepository
	marketID  int64
}

func NewPlayStoreService(
	ctx context.Context,
	credentialsJSON []byte,
	stateRepo repository.AppMarketStateRepository,
	marketID int64,
) (*PlayStoreService, error) {
	svc, err := androidpublisher.NewService(ctx, option.WithCredentialsJSON(credentialsJSON))
	if err != nil {
		return nil, fmt.Errorf("init play service: %w", err)
	}
	return &PlayStoreService{
		pub:       svc,
		edits:     svc.Edits,
		reviews:   svc.Reviews,
		tracks:    svc.Edits.Tracks,
		stateRepo: stateRepo,
		marketID:  marketID,
	}, nil
}

// getCurrentRelease — чистый запрос в Google Play.
func (p *PlayStoreService) getCurrentRelease(
	ctx context.Context,
	pkg string,
) (*domain.ReleaseInfo, error) {

	ed, err := p.edits.Insert(pkg, nil).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("play edit insert: %w", err)
	}
	defer func() {
		_ = p.edits.Delete(pkg, ed.Id).Context(ctx).Do()
	}()

	tr, err := p.tracks.Get(pkg, ed.Id, "production").Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("play tracks.get: %w", err)
	}
	if len(tr.Releases) == 0 {
		return nil, fmt.Errorf("play: no production release")
	}
	rel := tr.Releases[0]

	var maxCode int64
	for _, v := range rel.VersionCodes {
		if v > maxCode {
			maxCode = v
		}
	}

	m := playNameRe.FindStringSubmatch(rel.Name)
	var sem string
	if len(m) == 3 {
		sem = m[2]
	}

	return &domain.ReleaseInfo{
		Store:    domain.StoreGooglePlay,
		AppID:    pkg,
		Semantic: sem,
		Code:     maxCode,
	}, nil
}

// внутри PlayStoreService

// CheckNewRelease — возвращает релиз только если он НОВЕЕ, чем в app_market_states.
// Иначе (нет новой версии) → (nil, nil).
func (p *PlayStoreService) CheckNewRelease(
	ctx context.Context,
	app domain.App,
	packageName string,
	stateRepo repository.AppMarketStateRepository,
	marketID int64,
) (*domain.ReleaseInfo, error) {

	// читаем текущее состояние
	state, err := stateRepo.Get(ctx, app.ID, marketID)
	if err != nil {
		return nil, fmt.Errorf("get app_market_state: %w", err)
	}

	// текущая инфа из стора
	cur, err := p.getCurrentRelease(ctx, packageName)
	if err != nil {
		return nil, err
	}

	var lastCode int64
	if state != nil && state.LastVersion != nil {
		lastCode = *state.LastVersion
	}

	// если версия не выросла — ничего не возвращаем
	if cur.Code <= lastCode {
		return nil, nil
	}

	// есть новая версия — обновляем состояние
	if state == nil {
		state = &domain.AppMarketState{
			AppID:    app.ID,
			MarketID: marketID,
		}
	}
	state.LastVersion = &cur.Code

	if _, err := stateRepo.Save(ctx, state); err != nil {
		return nil, fmt.Errorf("save app_market_state: %w", err)
	}

	return cur, nil
}

// listReviewsRaw — твой старый ListReviews с параметром lastID.
func (p *PlayStoreService) listReviewsRaw(
	ctx context.Context,
	pkg, lastID string,
) ([]domain.Review, error) {

	call := p.reviews.List(pkg)
	resp, err := call.Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("play reviews.list: %w", err)
	}

	out := make([]domain.Review, 0, len(resp.Reviews))
stopLoop:
	for _, r := range resp.Reviews {
		if r.ReviewId == lastID {
			break stopLoop
		}
		if len(r.Comments) == 0 || r.Comments[0].UserComment == nil {
			continue
		}
		c := r.Comments[0].UserComment
		out = append(out, domain.Review{
			Store:        domain.StoreGooglePlay,
			AppID:        pkg,
			ReviewID:     r.ReviewId,
			AuthorName:   r.AuthorName,
			Rating:       int(c.StarRating),
			Text:         c.Text,
			LastModified: time.Unix(c.LastModified.Seconds, 0),
		})
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].LastModified.Before(out[j].LastModified)
	})
	return out, nil
}

func (p *PlayStoreService) FetchNewReviews(
	ctx context.Context,
	app domain.App,
	packageName string,
) ([]domain.Review, error) {

	state, err := p.stateRepo.Get(ctx, app.ID, p.marketID)
	if err != nil {
		return nil, fmt.Errorf("get app_market_state: %w", err)
	}

	var lastSeen string
	if state != nil && state.LastReviewID != nil {
		lastSeen = *state.LastReviewID
	}

	revs, err := p.listReviewsRaw(ctx, packageName, lastSeen)
	if err != nil {
		return nil, err
	}
	if len(revs) == 0 {
		return nil, nil
	}

	// Google обычно возвращает от новых к старым, но мы уже отсортировали по времени.
	last := revs[len(revs)-1]
	lastID := last.ReviewID

	if state == nil {
		state = &domain.AppMarketState{
			AppID:    app.ID,
			MarketID: p.marketID,
		}
	}
	state.LastReviewID = &lastID

	if _, err := p.stateRepo.Save(ctx, state); err != nil {
		return nil, fmt.Errorf("save app_market_state: %w", err)
	}

	return revs, nil
}
