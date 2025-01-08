package handlers

import "net/http"

func getTokenFromCookie(request *http.Request) (string, error) {
	cookie, err := request.Cookie("oauth2_token")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
