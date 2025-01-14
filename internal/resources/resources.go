package resources

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/antgobar/famcal/config"
)

func Load(router *http.ServeMux, config *config.Config) {
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveIndex(w, r)
	})
	router.HandleFunc("GET /privacy", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/html/privacy.html")
	})
	router.HandleFunc("GET /terms", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/html/terms.html")
	})
	router.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/img/favicon.ico")
	})
	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./static/html/index.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	host := r.Host

	baseURL := fmt.Sprintf("%s://%s", scheme, host)
	data := struct {
		APIBaseURL string
	}{
		APIBaseURL: baseURL,
	}

	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
