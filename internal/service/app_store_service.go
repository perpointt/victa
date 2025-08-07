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

	"github.com/golang-jwt/jwt/v5"
)

// AppStoreConfig настраивает AppStoreService.
type AppStoreConfig struct {
	KeyID      string // 10-символьный, начинается с "2"
	IssuerID   string // GUID из App Store Connect
	PrivatePEM []byte // содержимое .p8
	BaseURL    string
	HTTPClient *http.Client // если nil — будет создан с таймаутом 10s
}

// AppStoreService реализует StoreService для App Store.
type AppStoreService struct {
	cfg     AppStoreConfig
	signKey *ecdsa.PrivateKey
}

// NewAppStoreService создаёт AppStoreService.
func NewAppStoreService(cfg AppStoreConfig) (*AppStoreService, error) {
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
	return &AppStoreService{cfg: cfg, signKey: ecdsaKey}, nil
}

// jwtToken генерирует Bearer-токен для запросов.
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

// GetRelease возвращает последний production-релиз из App Store.
func (s *AppStoreService) GetRelease(
	ctx context.Context, appID string,
) (*domain.ReleaseInfo, error) {
	// 1. подготовить запрос
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

	// 2. выполнить
	resp, err := s.cfg.HTTPClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("asc %d: %s", resp.StatusCode, body)
	}

	// 3. разобрать JSON
	var raw struct {
		Data []struct {
			Attributes struct {
				VersionString string `json:"versionString"`
			}
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

	// 4. найти build в included
	var code int64
	for _, inc := range raw.Included {
		if inc.ID == buildID && inc.Type == "appStoreVersionBuilds" || inc.Type == "builds" {
			code, _ = strconv.ParseInt(inc.Attributes.Version, 10, 64)
			break
		}
	}

	// 5. вернуть unified модель
	return &domain.ReleaseInfo{
		Store:    domain.StoreAppStore,
		AppID:    appID,
		BundleID: "", // здесь можно доп. запросом получить bundle ID
		Semantic: ver,
		Code:     code,
	}, nil
}

// ListReviews возвращает отзывы из App Store.
func (s *AppStoreService) ListReviews(
	ctx context.Context, appID, lastSeenID string,
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
	defer func() {
		_ = resp.Body.Close()
	}()

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

	// от старых к новым
	sort.Slice(out, func(i, j int) bool {
		return out[i].LastModified.Before(out[j].LastModified)
	})
	return out, nil
}
