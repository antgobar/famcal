package resources

import (
	"io/fs"
	"log"
	"net/http"
)

func Load(router *http.ServeMux, assets fs.FS) {
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveFile(w, assets, "html/index.html")
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
