package main

import (
	"net/http"

	"nerdlp/engine.go"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	server := engine.New().Attach(r)

	http.ListenAndServe(":3000", server)
}
