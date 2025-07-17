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
	"strings"
	"time"

	"victa/internal/domain"

	"github.com/golang-jwt/jwt/v5"
)

type AppStoreConfig struct {
	KeyID      string // 10‑символьный, начинается с "2"
	IssuerID   string // GUID от App Store Connect
	PrivatePEM []byte // содержимое AuthKey_XXXX.p8
	BaseURL    string
	HTTP       *http.Client // необязательно; если nil — &http.Client{Timeout:10s}
}

type AppStoreService struct {
	cfg     AppStoreConfig
	signKey *ecdsa.PrivateKey
}

func NewAppStoreService(cfg AppStoreConfig) (*AppStoreService, error) {
	if cfg.HTTP == nil {
		cfg.HTTP = &http.Client{Timeout: 10 * time.Second}
	}
	block, _ := pem.Decode(cfg.PrivatePEM)
	k, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse p8: %w", err)
	}
	return &AppStoreService{
		cfg:     cfg,
		signKey: k.(*ecdsa.PrivateKey),
	}, nil
}

func (a *AppStoreService) jwtToken() (string, error) {
	claims := jwt.MapClaims{
		"iss": a.cfg.IssuerID,
		"exp": time.Now().Add(15 * time.Minute).Unix(),
		"aud": "appstoreconnect-v1",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["kid"] = a.cfg.KeyID
	return token.SignedString(a.signKey)
}

func (a *AppStoreService) GetLatestRelease(
	ctx context.Context,
	appID string, // 10‑значный numeric
) (*domain.AppStoreReleaseInfo, error) {

	j, _ := a.jwtToken()
	url := fmt.Sprintf(
		"%s/v1/apps/%s/appStoreVersions?filter[platform]=IOS&limit=20&include=build",
		a.cfg.BaseURL, appID)

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	req.Header.Set("Authorization", "Bearer "+j)
	resp, err := a.cfg.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("asc %d: %s", resp.StatusCode, b)
	}

	var raw struct {
		Data []struct {
			ID         string `json:"id"`
			Attributes struct {
				VersionString string `json:"versionString"` // "1.2.0"
			} `json:"attributes"`
			Relationships struct {
				Build struct { // <‑‑ singular
					Data struct {
						ID string `json:"id"`
					} `json:"data"`
				} `json:"build"`
			} `json:"relationships"`
		} `json:"data"`
		Included []struct {
			ID         string `json:"id"`
			Type       string `json:"type"` // "builds"
			Attributes struct {
				Version string `json:"version"` // build number "30"
			} `json:"attributes"`
		} `json:"included"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}
	if len(raw.Data) == 0 {
		return nil, fmt.Errorf("no appStoreVersions")
	}

	// берём последнюю по VersionString
	latest := raw.Data[0]
	for _, d := range raw.Data[1:] {
		if versionLT(latest.Attributes.VersionString, d.Attributes.VersionString) {
			latest = d
		}
	}
	semVer := latest.Attributes.VersionString

	// ищем build в included
	buildID := latest.Relationships.Build.Data.ID
	var buildNumStr string
	for _, inc := range raw.Included {
		if inc.ID == buildID && inc.Type == "builds" {
			buildNumStr = inc.Attributes.Version
			break
		}
	}
	buildNum, _ := strconv.ParseInt(buildNumStr, 10, 64)

	return &domain.AppStoreReleaseInfo{
		AppID: appID,
		Version: domain.PlayVersion{
			Semantic: semVer,
			Code:     buildNum,
		},
	}, nil
}

// versionLT reports whether a < b for strings like "1.2.3".
// Если какая‑то часть отсутствует, считаем её 0: "1.2" == "1.2.0".
func versionLT(a, b string) bool {
	pa := strings.Split(a, ".")
	pb := strings.Split(b, ".")
	max := len(pa)
	if len(pb) > max {
		max = len(pb)
	}

	for i := 0; i < max; i++ {
		var ai, bi int64
		if i < len(pa) {
			ai, _ = strconv.ParseInt(pa[i], 10, 64)
		}
		if i < len(pb) {
			bi, _ = strconv.ParseInt(pb[i], 10, 64)
		}
		if ai < bi {
			return true
		}
		if ai > bi {
			return false
		}
	}
	return false // равны
}

// -----------------------------------------------------------------------------
// Reviews
// -----------------------------------------------------------------------------

func (a *AppStoreService) ListReviewsSince(
	ctx context.Context,
	appID, lastSeenID string, // "" => всё, иначе отфильтруем
) ([]domain.PlayReview, error) {

	limit := 200
	j, _ := a.jwtToken()
	url := fmt.Sprintf("%s/v1/apps/%s/customerReviews?sort=-createdDate&limit=%d",
		a.cfg.BaseURL, appID, limit)

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	req.Header.Set("Authorization", "Bearer "+j)
	resp, err := a.cfg.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("asc %d: %s", resp.StatusCode, b)
	}

	var raw struct {
		Data []struct {
			ID         string `json:"id"`
			Attributes struct {
				Rating       int64  `json:"rating"`
				UserName     string `json:"userName"`
				Body         string `json:"body"`
				LastModified string `json:"lastModifiedDate"` // ISO8601
			} `json:"attributes"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	out := make([]domain.PlayReview, 0, len(raw.Data))
	for _, r := range raw.Data {
		if r.ID == lastSeenID {
			break
		}
		t, _ := time.Parse(time.RFC3339, r.Attributes.LastModified)
		out = append(out, domain.PlayReview{
			ReviewID:     r.ID,
			AuthorName:   r.Attributes.UserName,
			Rating:       r.Attributes.Rating,
			Text:         r.Attributes.Body,
			LastModified: t,
		})
	}

	// сортируем от старых к новым
	sort.Slice(out, func(i, j int) bool {
		return out[i].LastModified.Before(out[j].LastModified)
	})
	return out, nil
}
