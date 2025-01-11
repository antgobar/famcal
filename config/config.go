package config

import (
	"errors"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

type Config struct {
	GoogleOauth2Config *oauth2.Config
}

func LoadConfig() (*Config, error) {
	googleOauth2Config, err := loadGoogleOauth2Config()
	if err != nil {
		return nil, err
	}
	return &Config{GoogleOauth2Config: googleOauth2Config}, nil
}

func loadGoogleOauth2Config() (*oauth2.Config, error) {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURI := os.Getenv("GOOGLE_REDIRECT_URI")

	if clientID == "" {
		return nil, errors.New("environment variable GOOGLE_CLIENT_ID is required")
	}
	if clientSecret == "" {
		return nil, errors.New("environment variable GOOGLE_CLIENT_SECRET is required")
	}
	if redirectURI == "" {
		return nil, errors.New("environment variable GOOGLE_REDIRECT_URI is required")
	}

	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURI,
		Scopes:       []string{calendar.CalendarScope},
		Endpoint:     google.Endpoint,
	}, nil
}
