package middleware

import (
	"context"
	"net/http"
	"strings"

	"tech-ip-sem2/services/tasks/internal/client/authclient"
)

type AuthMiddleware struct {
	authClient *authclient.AuthClient
}

func NewAuthMiddleware(authClient *authclient.AuthClient) *AuthMiddleware {
	return &AuthMiddleware{
		authClient: authClient,
	}
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		ctx := r.Context()

		resp, statusCode, err := m.authClient.VerifyTokenWithHeader(ctx, token)
		if err != nil {
			http.Error(w, `{"error":"auth service unavailable"}`, http.StatusServiceUnavailable)
			return
		}

		if statusCode != http.StatusOK || !resp.Valid {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(ctx, "subject", resp.Subject)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
