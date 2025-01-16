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
	Env                string
	ServerAddr         string
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
	env, err := loadSimpleStringEnvVar("ENV")
	if err != nil {
		return nil, err
	}
	return &Config{GoogleOauth2Config: googleOauth2Config, Env: *env, ServerAddr: setServerAddr(*env)}, nil
}

func loadGoogleOauth2Config() (*oauth2.Config, error) {
	clientID, err := loadSimpleStringEnvVar("GOOGLE_CLIENT_ID")
	if err != nil {
		return nil, err
	}
	clientSecret, err := loadSimpleStringEnvVar("GOOGLE_CLIENT_SECRET")
	if err != nil {
		return nil, err
	}
	redirectURI, err := loadSimpleStringEnvVar("GOOGLE_REDIRECT_URI")
	if err != nil {
		return nil, err
	}

	return &oauth2.Config{
		ClientID:     *clientID,
		ClientSecret: *clientSecret,
		RedirectURL:  *redirectURI,
		Scopes:       []string{calendar.CalendarScope, "profile", "email"},
		Endpoint:     google.Endpoint,
	}, nil
}

func newEnvVarErr(varName string) string {
	return "environment variable " + varName + " is required"
}

func loadSimpleStringEnvVar(name string) (*string, error) {
	envVar := os.Getenv(name)
	if envVar == "" {
		return nil, errors.New(newEnvVarErr(name))
	}
	return &envVar, nil
}

func setServerAddr(env string) string {
	if env == "development" {
		return "localhost:8090"
	}
	return "0.0.0.0:8090"
}
