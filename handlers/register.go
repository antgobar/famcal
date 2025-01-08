package handlers

import "net/http"

func Register(router *http.ServeMux) {
	router.HandleFunc("GET /members", getMembers)
	router.HandleFunc("POST /members", addMember)
	router.HandleFunc("GET /calendars", getCalendars)
	router.HandleFunc("GET /auth/google/callback", handleGoogleAuthCallback)
	router.HandleFunc("GET /auth/google/auth", authHandler)
}
