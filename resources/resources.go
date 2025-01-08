package resources

import "net/http"

func Load(router *http.ServeMux) {
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/html/index.html")
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
