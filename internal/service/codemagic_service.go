package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"victa/internal/domain"
)

type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type CodemagicService struct {
	client  HTTPDoer
	baseURL string
	ttl     time.Duration
}

func NewCodemagicService(baseURL string) *CodemagicService {
	return &CodemagicService{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: strings.TrimRight(baseURL, "/"),
		ttl:     7 * 24 * time.Hour,
	}
}

func (s *CodemagicService) WithHTTPClient(c HTTPDoer) *CodemagicService {
	s.client = c
	return s
}

func (s *CodemagicService) WithArtifactTTL(d time.Duration) *CodemagicService {
	s.ttl = d
	return s
}

func (s *CodemagicService) GetBuildByID(
	ctx context.Context,
	buildID, apiKey string,
) (*domain.CodemagicBuildResponse, error) {

	url := fmt.Sprintf("%s/builds/%s", s.baseURL, buildID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("x-auth-token", apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("codemagic %d: %s", resp.StatusCode, body)
	}

	var out domain.CodemagicBuildResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &out, nil
}

func (s *CodemagicService) GetArtifactPublicURL(
	ctx context.Context,
	path, apiKey string,
) (string, error) {

	expires := time.Now().Add(s.ttl).Unix()
	payload, _ := json.Marshal(struct {
		ExpiresAt int64 `json:"expiresAt"`
	}{expires})

	url := fmt.Sprintf("%s/artifacts/%s/public-url", s.baseURL, strings.TrimPrefix(path, "/"))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("x-auth-token", apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("execute request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("codemagic %d: %s", resp.StatusCode, body)
	}

	var out struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}
	return out.URL, nil
}
