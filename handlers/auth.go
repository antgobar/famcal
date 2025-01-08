package handlers

import (
	"fmt"
	"net/http"

	"github.com/antgobar/famcal/integrations"
)

func handleGoogleAuthCallback(w http.ResponseWriter, r *http.Request) {
	authCode := r.URL.Query().Get("code")
	if authCode == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	token, err := integrations.HandleRequestToken(authCode)
	if err != nil {
		http.Error(w, "Error retrieving token", http.StatusBadRequest)
		return
	}

	err = integrations.SaveToken("token.json", token)
	if err != nil {
		http.Error(w, "Error caching token", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	url, err := integrations.GoogleAuthUrl()
	if err != nil || url == "" {
		http.Error(w, fmt.Sprintf("Auth error: %v", err), http.StatusUnauthorized)
		return
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
