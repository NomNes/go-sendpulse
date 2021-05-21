package go_sendpulse

import (
	"context"

	"github.com/go-resty/resty/v2"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	tokenURL = "https://api.sendpulse.com/oauth/access_token"
	ApiURL   = "https://api.sendpulse.com"
)

type Client struct {
	*resty.Client
}

func New(ctx context.Context, clientId, clientSecret string) Client {
	config := clientcredentials.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		TokenURL:     tokenURL,
	}
	client := resty.NewWithClient(config.Client(ctx))
	client.SetHostURL(ApiURL)
	return Client{client}
}
