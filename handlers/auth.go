package handlers

import (
	"fmt"
	"net/http"

	"github.com/antgobar/famcal/provider"
)

func (h Handler) handleGoogleAuthCallback(w http.ResponseWriter, r *http.Request) {
	authCode := r.URL.Query().Get("code")
	if authCode == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	token, err := provider.HandleCallback(authCode, h.config.GoogleOauth2Config)
	if err != nil {
		http.Error(w, "Error retrieving token", http.StatusBadRequest)
		return
	}
	SetOauth2TokenCookie(w, token)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h Handler) googleAuthHandler(w http.ResponseWriter, r *http.Request) {
	url, err := provider.AuthUrl(*h.config.GoogleOauth2Config)
	if err != nil || url == "" {
		http.Error(w, fmt.Sprintf("Auth error: %v", err), http.StatusUnauthorized)
		return
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h Handler) logOutGoogle(w http.ResponseWriter, r *http.Request) {
	ClearOauth2TokenCookie(w)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
