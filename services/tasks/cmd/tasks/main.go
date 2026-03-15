package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"tech-ip-sem2/services/tasks/internal/client/authclient"
	"tech-ip-sem2/services/tasks/internal/handler"
	"tech-ip-sem2/services/tasks/internal/middleware"
	"tech-ip-sem2/services/tasks/internal/storage"
	sharedMiddleware "tech-ip-sem2/shared/middleware"
)

func main() {
	port := os.Getenv("TASKS_PORT")
	if port == "" {
		port = "8082"
	}

	authBaseURL := os.Getenv("AUTH_BASE_URL")
	if authBaseURL == "" {
		authBaseURL = "http://localhost:8081"
	}

	taskStorage := storage.NewMemoryStorage()
	authClient := authclient.NewAuthClient(authBaseURL)

	tasksHandler := handler.NewTasksHandler(taskStorage)

	authMiddleware := middleware.NewAuthMiddleware(authClient)

	mux := http.NewServeMux()

	mux.Handle("POST /v1/tasks", authMiddleware.RequireAuth(http.HandlerFunc(tasksHandler.CreateTask)))
	mux.Handle("GET /v1/tasks", authMiddleware.RequireAuth(http.HandlerFunc(tasksHandler.GetTasks)))
	mux.Handle("GET /v1/tasks/{id}", authMiddleware.RequireAuth(http.HandlerFunc(tasksHandler.GetTask)))
	mux.Handle("PATCH /v1/tasks/{id}", authMiddleware.RequireAuth(http.HandlerFunc(tasksHandler.UpdateTask)))
	mux.Handle("DELETE /v1/tasks/{id}", authMiddleware.RequireAuth(http.HandlerFunc(tasksHandler.DeleteTask)))

	handler := sharedMiddleware.RequestIDMiddleware(mux)

	fmt.Printf("Tasks service running on port %s, using auth at %s\n", port, authBaseURL)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
