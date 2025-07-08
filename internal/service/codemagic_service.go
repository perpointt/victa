package service

import (
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
