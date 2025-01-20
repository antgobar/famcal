package config

import (
	"errors"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/people/v1"
)

type Config struct {
	GoogleOauth2Config *oauth2.Config
	Env                string
	ServerAddr         string
	Host               string
}

func MustLoadConfig() *Config {
	googleOauth2Config := mustLoadGoogleOauth2Config()
	env := mustLoadEnvVar("ENV")
	return &Config{GoogleOauth2Config: googleOauth2Config, Env: *env, ServerAddr: setServerAddr(*env)}
}

func mustLoadEnvVar(varName string) *string {
	env, err := loadSimpleStringEnvVar(varName)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	return env
}

func mustLoadGoogleOauth2Config() *oauth2.Config {
	clientID := mustLoadEnvVar("GOOGLE_CLIENT_ID")
	clientSecret := mustLoadEnvVar("GOOGLE_CLIENT_SECRET")
	redirectURI := mustLoadEnvVar("GOOGLE_REDIRECT_URI")
	return &oauth2.Config{
		ClientID:     *clientID,
		ClientSecret: *clientSecret,
		RedirectURL:  *redirectURI,
		Scopes: []string{
			calendar.CalendarScope,
			people.UserEmailsReadScope,
			people.UserinfoProfileScope,
		},
		Endpoint: google.Endpoint,
	}
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
	switch env {
	case "production":
		return "0.0.0.0:8090"
	default:
		return "localhost:8090"
	}
}
