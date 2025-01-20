package resources

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"github.com/antgobar/famcal/config"
)

func Load(router *http.ServeMux, config *config.Config, assets fs.FS) {
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveIndex(w, r, config, assets)
	})
	router.HandleFunc("GET /privacy", func(w http.ResponseWriter, r *http.Request) {
		serveFile(w, assets, "html/privacy.html")
	})
	router.HandleFunc("GET /terms", func(w http.ResponseWriter, r *http.Request) {
		serveFile(w, assets, "html/terms.html")
	})
	router.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		serveFile(w, assets, "img/favicon.ico")
	})
	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(assets))))
}

func generateBaseUrl(r *http.Request, cfg *config.Config) string {
	scheme := "http"
	if env := cfg.Env; env == "production" {
		scheme = "https"
	}
	host := r.Host

	return fmt.Sprintf("%s://%s", scheme, host)
}

func serveIndex(w http.ResponseWriter, r *http.Request, cfg *config.Config, assets fs.FS) {
	tmpl, err := template.ParseFS(assets, "html/index.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	baseURL := generateBaseUrl(r, cfg)
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

func serveFile(w http.ResponseWriter, assets fs.FS, path string) {
	content, err := fs.ReadFile(assets, path)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "File not found: "+path, http.StatusNotFound)
		return
	}
	if ext := http.DetectContentType(content); ext != "" {
		w.Header().Set("Content-Type", ext)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}
