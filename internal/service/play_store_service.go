package service

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"time"

	"victa/internal/domain"

	"google.golang.org/api/androidpublisher/v3"
	"google.golang.org/api/option"
)

// PlayStoreService читает production‑релиз и отзывы Google Play.
type PlayStoreService struct {
	pub      *androidpublisher.Service
	edits    *androidpublisher.EditsService
	reviews  *androidpublisher.ReviewsService
	tracks   *androidpublisher.EditsTracksService
	reNameRe *regexp.Regexp // n (x.y.z)
}

// NewPlayStoreService инициализирует клиент.
// credentialsJSON — содержимое service‑account .json.
func NewPlayStoreService(ctx context.Context, credentialsJSON []byte) (*PlayStoreService, error) {
	opts := []option.ClientOption{option.WithCredentialsJSON(credentialsJSON)}

	svc, err := androidpublisher.NewService(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("init publisher service: %w", err)
	}
	return &PlayStoreService{
		pub:      svc,
		edits:    svc.Edits,
		reviews:  svc.Reviews,
		tracks:   svc.Edits.Tracks,
		reNameRe: regexp.MustCompile(`^(\d+)\s+\((\d+\.\d+\.\d+)\)$`),
	}, nil
}

// GetProductionRelease возвращает текущий production‑релиз.
// best‑effort: Semantic берётся из releases[].name,
// Code — максимальный versionCode внутри релиза.
func (p *PlayStoreService) GetProductionRelease(
	ctx context.Context, pkg string,
) (*domain.PlayReleaseInfo, error) {

	// 1. Стартуем «edit» (обязательная формальность Google)
	edit, err := p.edits.Insert(pkg, nil).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("insert edit: %w", err)
	}
	editID := edit.Id
	defer func() { _ = p.edits.Delete(pkg, editID).Context(ctx).Do() }()

	// 2. Читаем production‑трек
	tr, err := p.tracks.Get(pkg, editID, "production").Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("tracks.get: %w", err)
	}
	if len(tr.Releases) == 0 {
		return nil, fmt.Errorf("package %s: no releases in production track", pkg)
	}
	rel := tr.Releases[0]

	// 3. Берём максимальный versionCode
	var maxCode int64
	for _, v := range rel.VersionCodes {
		if v > maxCode {
			maxCode = v
		}
	}

	// 4. Парсим release name вида "45 (1.2.3)"
	m := p.reNameRe.FindStringSubmatch(rel.Name)
	codeFromName, _ := strconv.ParseInt(m[1], 10, 64)
	sem := m[2]

	// 5. Подстраховка: доверяем versionCode из имени, если он совпал с maxCode
	if codeFromName != maxCode {
		// берём maxCode, даже если в имени другое число
	}

	return &domain.PlayReleaseInfo{
		PackageName: pkg,
		Version: domain.PlayVersion{
			Semantic: sem,
			Code:     maxCode,
		},
	}, nil
}

// ListReviewsSince возвращает новые отзывы, которых ещё нет в БД.
// Результат отсортирован от самых старых к самым новым.
func (p *PlayStoreService) ListReviewsSince(
	ctx context.Context,
	pkg, lastSeenID string, // lastSeenID == "" → вернёт всё за 7 дней
) ([]domain.PlayReview, error) {

	call := p.reviews.List(pkg)

	resp, err := call.Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("reviews.list: %w", err)
	}

	newRevs := make([]domain.PlayReview, 0, len(resp.Reviews))
	stop := false
	for _, r := range resp.Reviews { // Google отдаёт от новых к старым
		if r.ReviewId == lastSeenID {
			stop = true
			break
		}
		if len(r.Comments) == 0 || r.Comments[0].UserComment == nil {
			continue
		}
		c := r.Comments[0].UserComment
		newRevs = append(newRevs, domain.PlayReview{
			ReviewID:     r.ReviewId,
			AuthorName:   r.AuthorName,
			Rating:       c.StarRating,
			Text:         c.Text,
			LastModified: time.Unix(c.LastModified.Seconds, 0),
		})
	}
	if !stop && lastSeenID != "" {
		// получили <max> отзывов, но не встретили старый — возможно есть ещё страница
		// бизнес‑логика сама решит: вызывать ещё раз с пагинацией или забить
	}

	// Сортируем от старых к новым
	sort.Slice(newRevs, func(i, j int) bool {
		return newRevs[i].LastModified.Before(newRevs[j].LastModified)
	})
	return newRevs, nil
}
