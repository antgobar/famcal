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
	ServerAddress      string
}

func Load() (*Config, error) {
	googleOauth2Config, err := loadGoogleOauth2Config()
	if err != nil {
		return nil, err
	}
	serverAddr, err := loadServerConfig()
	if err != nil {
		return nil, err
	}
	return &Config{GoogleOauth2Config: googleOauth2Config, ServerAddress: *serverAddr}, nil
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

func newEnvVarErr(varName string) string {
	return "environment variable " + varName + " is required"
}

func loadServerConfig() (*string, error) {
	serverAddr := os.Getenv("SERVER_ADDR")
	if serverAddr == "" {
		return nil, errors.New(newEnvVarErr("SERVER_ADDR"))
	}
	return &serverAddr, nil
}
