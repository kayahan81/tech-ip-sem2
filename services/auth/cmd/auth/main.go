package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"tech-ip-sem2/services/auth/internal/handler"
	sharedMiddleware "tech-ip-sem2/shared/middleware"
)

func main() {
	port := os.Getenv("AUTH_PORT")
	if port == "" {
		port = "8081"
	}

	mux := http.NewServeMux()

	// Хендлеры
	authHandler := handler.NewAuthHandler()
	mux.HandleFunc("POST /v1/auth/login", authHandler.Login)
	mux.HandleFunc("GET /v1/auth/verify", authHandler.Verify)

	// Глобальные middleware
	handler := sharedMiddleware.RequestIDMiddleware(mux)

	fmt.Printf("Auth service running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
