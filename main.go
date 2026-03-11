package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var users = []User{
	{ID: 1, Name: "Nome 1"},
	{ID: 2, Name: "Nome 2"},
}

func main() {

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Teste"))
	})

	r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(users)
	})

	http.ListenAndServe(":3000", r)
}
