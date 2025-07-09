package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"victa/internal/domain"
)

// CodemagicService умеет запрашивать данные о сборке по её ID.
type CodemagicService struct {
	client  *http.Client
	baseURL string
}

// NewCodemagicService создаёт сервис с указанным API-токеном.
// baseURL обычно "https://api.codemagic.io".
func NewCodemagicService(baseURL string) *CodemagicService {
	return &CodemagicService{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: strings.TrimRight(baseURL, "/"),
	}
}

// GetBuildByID запрашивает GET /builds/:id и разбирает JSON в CodemagicBuildResponse.
func (s *CodemagicService) GetBuildByID(buildID, apiKey string) (*domain.CodemagicBuildResponse, error) {
	url := fmt.Sprintf("%s/builds/%s", s.baseURL, buildID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	// Codemagic ждёт заголовок x-auth-token
	req.Header.Set("x-auth-token", apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("codemagic API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("codemagic API returned %d: %s", resp.StatusCode, string(body))
	}

	var out domain.CodemagicBuildResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("failed to decode codemagic response: %w", err)
	}
	return &out, nil
}

func (s *CodemagicService) GetArtifactPublicURL(path, apiKey string) (string, error) {
	ttl := 7 * 24 * time.Hour
	expires := time.Now().Add(ttl).Unix()

	payload, _ := json.Marshal(struct {
		ExpiresAt int64 `json:"expiresAt"`
	}{expires})

	url := fmt.Sprintf("%s/artifacts/%s/public-url", s.baseURL, strings.TrimPrefix(path, "/"))

	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("x-auth-token", apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("codemagic request: %w", err)
	}
	defer resp.Body.Close()

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
