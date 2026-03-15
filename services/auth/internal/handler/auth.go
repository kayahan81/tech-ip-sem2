package handler

import (
	"encoding/json"
	"net/http"
	"strings"
)

const validToken = "demo-token"

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type VerifyResponse struct {
	Valid   bool   `json:"valid"`
	Subject string `json:"subject,omitempty"`
	Error   string `json:"error,omitempty"`
}

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp := LoginResponse{
		AccessToken: validToken,
		TokenType:   "Bearer",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) Verify(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")

	resp := VerifyResponse{}
	status := http.StatusOK

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		resp.Valid = false
		resp.Error = "unauthorized"
		status = http.StatusUnauthorized
	} else {
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == validToken {
			resp.Valid = true
			resp.Subject = "student"
		} else {
			resp.Valid = false
			resp.Error = "unauthorized"
			status = http.StatusUnauthorized
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}
