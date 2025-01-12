package handlers

import (
	"errors"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

const oauth2CookieName string = "oauth2_token"

func GetOauth2CookieValue(request *http.Request) (*string, error) {
	cookie, err := request.Cookie(oauth2CookieName)
	if err != nil {
		return nil, errors.New(err.Error() + oauth2CookieName)
	}
	return &cookie.Value, nil
}

func SetOauth2TokenCookie(w http.ResponseWriter, token *oauth2.Token) {
	cookie := &http.Cookie{
		Name:     oauth2CookieName,
		Value:    token.AccessToken,
		Expires:  time.Now().Add(time.Duration(time.Until(token.Expiry))),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	}
	http.SetCookie(w, cookie)
}
