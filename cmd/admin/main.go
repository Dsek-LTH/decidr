package main

import (
	"log"
	"net/http"

	"github.com/Dsek-LTH/decidr/internal/routes"
	"github.com/Dsek-LTH/decidr/internal/templates"
)

func loggingMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf(
			"%s -> %s",
			r.RemoteAddr,
			r.URL.String(),
		)

		next.ServeHTTP(w, r)
	})
}

func main() {
	templateRenderer := templates.NewTemplateRenderer()

	mux := http.NewServeMux()

	routes.RegisterRoutes(mux, templateRenderer)

	fs := http.FileServer(http.Dir("./web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	logged := loggingMiddleWare(mux)

	log.Println("Server listening on :11337")
	log.Fatal(http.ListenAndServe(":11337", logged))
}
