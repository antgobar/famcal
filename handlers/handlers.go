package handlers

import (
	"net/http"

	"github.com/antgobar/famcal/config"
)

type Handler struct {
	config *config.Config
}

func newHandler(config *config.Config) Handler {
	return Handler{config: config}
}

func Register(mux *http.ServeMux, config *config.Config) {
	h := newHandler(config)
	mux.HandleFunc("GET /members", getMembers)
	mux.HandleFunc("POST /members", addMember)
	mux.HandleFunc("GET /calendars", h.getCalendars)
	mux.HandleFunc("GET /events", h.getEvents)
	mux.HandleFunc("GET /auth/google/callback", h.handleGoogleAuthCallback)
	mux.HandleFunc("GET /auth/google/auth", h.authHandler)
}
