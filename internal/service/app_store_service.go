package service

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"time"

	"victa/internal/domain"
	"victa/internal/repository"

	"github.com/golang-jwt/jwt/v5"
)

type AppStoreConfig struct {
	KeyID      string
	IssuerID   string
	PrivatePEM []byte
	BaseURL    string
	HTTPClient *http.Client
}

type AppStoreService struct {
	cfg       AppStoreConfig
	signKey   *ecdsa.PrivateKey
	stateRepo repository.AppMarketStateRepository
	marketID  int64
}

func NewAppStoreService(
	cfg AppStoreConfig,
	stateRepo repository.AppMarketStateRepository,
	marketID int64,
) (*AppStoreService, error) {
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = &http.Client{Timeout: 10 * time.Second}
	}
	block, _ := pem.Decode(cfg.PrivatePEM)
	if block == nil {
		return nil, fmt.Errorf("invalid PEM")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse p8: %w", err)
	}
	ecdsaKey, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("not ECDSA key")
	}
	return &AppStoreService{
		cfg:       cfg,
		signKey:   ecdsaKey,
		stateRepo: stateRepo,
		marketID:  marketID,
	}, nil
}

func (s *AppStoreService) jwtToken() (string, error) {
	claims := jwt.MapClaims{
		"iss": s.cfg.IssuerID,
		"exp": time.Now().Add(15 * time.Minute).Unix(),
		"aud": "appstoreconnect-v1",
	}
	t := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	t.Header["kid"] = s.cfg.KeyID
	return t.SignedString(s.signKey)
}

func (s *AppStoreService) getCurrentRelease(
	ctx context.Context,
	appID string,
) (*domain.ReleaseInfo, error) {
	token, err := s.jwtToken()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(
		"%s/v1/apps/%s/appStoreVersions?filter[platform]=IOS&limit=1&include=build",
		s.cfg.BaseURL, appID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := s.cfg.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("asc %d: %s", resp.StatusCode, body)
	}

	var raw struct {
		Data []struct {
			Attributes struct {
				VersionString string `json:"versionString"`
			} `json:"attributes"`
			Relationships struct {
				Build struct {
					Data struct {
						ID string `json:"id"`
					} `json:"data"`
				} `json:"build"`
			} `json:"relationships"`
		} `json:"data"`
		Included []struct {
			ID         string `json:"id"`
			Type       string `json:"type"`
			Attributes struct {
				Version string `json:"version"`
			} `json:"attributes"`
		} `json:"included"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}
	if len(raw.Data) == 0 {
		return nil, fmt.Errorf("no appStoreVersions")
	}

	ver := raw.Data[0].Attributes.VersionString
	buildID := raw.Data[0].Relationships.Build.Data.ID

	var code int64
	for _, inc := range raw.Included {
		if inc.ID == buildID && (inc.Type == "appStoreVersionBuilds" || inc.Type == "builds") {
			code, _ = strconv.ParseInt(inc.Attributes.Version, 10, 64)
			break
		}
	}

	return &domain.ReleaseInfo{
		Store:    domain.StoreAppStore,
		AppID:    appID,
		BundleID: "",
		Semantic: ver,
		Code:     code,
	}, nil
}

func (s *AppStoreService) CheckNewRelease(
	ctx context.Context,
	app domain.App,
	appID string,
	bundleID string,
	stateRepo repository.AppMarketStateRepository,
	marketID int64,
) (*domain.ReleaseInfo, error) {

	cur, err := s.getCurrentRelease(ctx, appID)
	if err != nil {
		return nil, err
	}

	state, err := stateRepo.Get(ctx, app.ID, marketID)
	if err != nil {
		return nil, fmt.Errorf("get app_market_state: %w", err)
	}

	var lastCode int64
	if state != nil && state.LastVersion != nil {
		lastCode = *state.LastVersion
	}

	if cur.Code <= lastCode {
		return nil, nil
	}

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

	// заодно проставим bundleID (если нужно)
	cur.BundleID = bundleID
	return cur, nil
}

func (s *AppStoreService) listReviewsRaw(
	ctx context.Context,
	appID, lastSeenID string,
) ([]domain.Review, error) {
	token, err := s.jwtToken()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf(
		"%s/v1/apps/%s/customerReviews?sort=-createdDate&limit=200",
		s.cfg.BaseURL, appID,
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := s.cfg.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("asc %d: %s", resp.StatusCode, body)
	}

	var raw struct {
		Data []struct {
			ID         string `json:"id"`
			Attributes struct {
				Rating           int64  `json:"rating"`
				UserName         string `json:"userName"`
				Body             string `json:"body"`
				LastModifiedDate string `json:"lastModifiedDate"`
			} `json:"attributes"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	out := make([]domain.Review, 0, len(raw.Data))
	for _, r := range raw.Data {
		if r.ID == lastSeenID {
			break
		}
		t, _ := time.Parse(time.RFC3339, r.Attributes.LastModifiedDate)
		out = append(out, domain.Review{
			Store:        domain.StoreAppStore,
			AppID:        appID,
			ReviewID:     r.ID,
			AuthorName:   r.Attributes.UserName,
			Rating:       int(r.Attributes.Rating),
			Text:         r.Attributes.Body,
			LastModified: t,
		})
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].LastModified.Before(out[j].LastModified)
	})
	return out, nil
}

func (s *AppStoreService) FetchNewReviews(
	ctx context.Context,
	app domain.App,
	appStoreID string,
) ([]domain.Review, error) {

	state, err := s.stateRepo.Get(ctx, app.ID, s.marketID)
	if err != nil {
		return nil, fmt.Errorf("get app_market_state: %w", err)
	}

	var lastSeen string
	if state != nil && state.LastReviewID != nil {
		lastSeen = *state.LastReviewID
	}

	revs, err := s.listReviewsRaw(ctx, appStoreID, lastSeen)
	if err != nil {
		return nil, err
	}
	if len(revs) == 0 {
		return nil, nil
	}

	// последний по времени отзыв
	last := revs[len(revs)-1]
	lastID := last.ReviewID // строка

	if state == nil {
		state = &domain.AppMarketState{
			AppID:    app.ID,
			MarketID: s.marketID,
		}
	}
	state.LastReviewID = &lastID

	if _, err := s.stateRepo.Save(ctx, state); err != nil {
		return nil, fmt.Errorf("save app_market_state: %w", err)
	}

	return revs, nil
}
