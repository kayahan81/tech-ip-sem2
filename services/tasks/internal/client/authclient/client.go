package authclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"tech-ip-sem2/shared/httpx"
)

type AuthClient struct {
	client  *httpx.Client
	baseURL string
}

type VerifyResponse struct {
	Valid   bool   `json:"valid"`
	Subject string `json:"subject,omitempty"`
	Error   string `json:"error,omitempty"`
}

func NewAuthClient(baseURL string) *AuthClient {
	return &AuthClient{
		client:  httpx.NewClient(baseURL, 3*time.Second),
		baseURL: baseURL,
	}
}

// VerifyToken проверяет токен через Auth service
func (c *AuthClient) VerifyToken(ctx context.Context, token string) (*VerifyResponse, int, error) {
	var resp VerifyResponse

	// Создаем запрос вручную, чтобы установить заголовок Authorization
	url := c.baseURL + "/v1/auth/verify"

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	// Прокидываем request-id из контекста
	if requestID, ok := ctx.Value("requestID").(string); ok && requestID != "" {
		req.Header.Set("X-Request-ID", requestID)
	}

	// Выполняем запрос - используем httpClient напрямую
	httpResp, err := c.client.Do(req) // Метод Do существует в httpx.Client
	if err != nil {
		return nil, 0, fmt.Errorf("do request: %w", err)
	}
	defer httpResp.Body.Close()

	// Декодируем ответ
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, httpResp.StatusCode, fmt.Errorf("decode response: %w", err)
	}

	return &resp, httpResp.StatusCode, nil
}

// VerifyTokenWithHeader - алиас для VerifyToken (для совместимости)
func (c *AuthClient) VerifyTokenWithHeader(ctx context.Context, token string) (*VerifyResponse, int, error) {
	return c.VerifyToken(ctx, token)
}
