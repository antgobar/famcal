package provider

import (
	"context"
	"errors"

	"golang.org/x/oauth2"
)

func HandleCallback(authCode string, config *oauth2.Config) (*oauth2.Token, error) {
	token, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, errors.New("unable to retrieve token from web")
	}
	return token, nil
}

func AuthUrl(config oauth2.Config) (string, error) {
	return config.AuthCodeURL("state-token", oauth2.AccessTypeOffline), nil
}
