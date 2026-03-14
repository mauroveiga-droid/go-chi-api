package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

var startTime = time.Now()

type EchoRequest struct {
	Payload string `json:"payload"`
}

type EchoResponse struct {
	Payload     string `json:"payload"`
	ProcessedAt string `json:"processed_at"`
}

func main() {
	r := chi.NewRouter()

	r.Use(loggingMiddleware)

	r.Get("/", helloWorld)
	r.Get("/health", healthCheck)
	r.Get("/hello/{name}", helloName)
	r.Get("/user/{id:[0-9]+}", userProfile)
	r.Get("/search", searchHandler)

	r.With(authMiddleware).Post("/echo", echoHandler)

	fmt.Println("Servidor a correr em http://localhost:8080")

	http.ListenAndServe(":8080", r)
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Olá Mundo!"))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(startTime)

	response := map[string]string{
		"status": "up",
		"uptime": uptime.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func helloName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	w.Write([]byte(fmt.Sprintf("Olá, %s!", name)))
}

func userProfile(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	w.Write([]byte(fmt.Sprintf("User Profile: %s", id)))
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	page := r.URL.Query().Get("page")

	if page == "" {
		page = "1"
	}

	response := fmt.Sprintf("Searching for '%s' on page %s", query, page)

	w.Write([]byte(response))
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	var req EchoRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Payload == "" {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	response := EchoResponse{
		Payload:     req.Payload,
		ProcessedAt: time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("X-App-Token")

		if token != "secret123" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Unauthorized: token inválido ou ausente",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(start)

		fmt.Printf("[%s] - %s - %v\n", r.Method, r.URL.Path, duration)
	})
}
