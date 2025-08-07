package service

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"time"

	"victa/internal/domain"

	"google.golang.org/api/androidpublisher/v3"
	"google.golang.org/api/option"
)

var playNameRe = regexp.MustCompile(`^(\d+)\s+\((\d+\.\d+\.\d+)\)$`)

type PlayStoreService struct {
	pub     *androidpublisher.Service
	edits   *androidpublisher.EditsService
	reviews *androidpublisher.ReviewsService
	tracks  *androidpublisher.EditsTracksService
}

// NewPlayStoreService возвращает клиент Google Play.
func NewPlayStoreService(ctx context.Context, credentialsJSON []byte) (*PlayStoreService, error) {
	svc, err := androidpublisher.NewService(ctx, option.WithCredentialsJSON(credentialsJSON))
	if err != nil {
		return nil, fmt.Errorf("init play service: %w", err)
	}
	return &PlayStoreService{
		pub:     svc,
		edits:   svc.Edits,
		reviews: svc.Reviews,
		tracks:  svc.Edits.Tracks,
	}, nil
}

// GetRelease реализует StoreService.GetRelease.
func (p *PlayStoreService) GetRelease(ctx context.Context, pkg string) (*domain.ReleaseInfo, error) {
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

	// max versionCode
	var maxCode int64
	for _, v := range rel.VersionCodes {
		if v > maxCode {
			maxCode = v
		}
	}

	// try parse name "123 (1.2.3)"
	m := playNameRe.FindStringSubmatch(rel.Name)
	var sem string
	if len(m) == 3 {
		sem = m[2]
	} else {
		sem = "" // не удалось
	}

	return &domain.ReleaseInfo{
		Store:    domain.StoreGooglePlay,
		AppID:    pkg,
		Semantic: sem,
		Code:     maxCode,
	}, nil
}

// ListReviews реализует StoreService.ListReviews.
func (p *PlayStoreService) ListReviews(ctx context.Context, pkg, lastID string) ([]domain.Review, error) {
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
	// от старых к новым
	sort.Slice(out, func(i, j int) bool {
		return out[i].LastModified.Before(out[j].LastModified)
	})
	return out, nil
}
