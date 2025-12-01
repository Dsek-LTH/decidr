package routes

import (
	"net/http"

	"github.com/Dsek-LTH/decidr/internal/templates"
)

func RegisterRoutes(mux *http.ServeMux, renderer *templates.TemplateRenderer) {
	mux.HandleFunc(
		"/",
		func(w http.ResponseWriter, req *http.Request) {
			renderer.Render(w, "home", nil)
		},
	)
}

