package config

import (
	"errors"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

type Config struct {
	GoogleOauth2Config *oauth2.Config
}

func MustLoadConfig() *Config {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Error creating config: %v", err)
	}
	return config
}

func loadConfig() (*Config, error) {
	googleOauth2Config, err := loadGoogleOauth2Config()
	if err != nil {
		return nil, err
	}
	return &Config{GoogleOauth2Config: googleOauth2Config}, nil
}

func loadGoogleOauth2Config() (*oauth2.Config, error) {
	clientID, err := loadSimpleStringEnvVar("GOOGLE_CLIENT_ID", false)
	if err != nil {
		return nil, err
	}
	clientSecret, err := loadSimpleStringEnvVar("GOOGLE_CLIENT_SECRET", false)
	if err != nil {
		return nil, err
	}
	redirectURI, err := loadSimpleStringEnvVar("GOOGLE_REDIRECT_URI", false)
	if err != nil {
		return nil, err
	}

	return &oauth2.Config{
		ClientID:     *clientID,
		ClientSecret: *clientSecret,
		RedirectURL:  *redirectURI,
		Scopes:       []string{calendar.CalendarScope},
		Endpoint:     google.Endpoint,
	}, nil
}

func newEnvVarErr(varName string) string {
	return "environment variable " + varName + " is required"
}

func loadSimpleStringEnvVar(name string, canBeEmpty bool) (*string, error) {
	envVar := os.Getenv(name)
	if envVar == "" && !canBeEmpty {
		return nil, errors.New(newEnvVarErr(name))
	}
	return &envVar, nil
}
