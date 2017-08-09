package cf

import (
	"context"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func ClientCredentials(clientId, secret string) option {
	return option{
		apply: func(client *Client) error {
			conf := &clientcredentials.Config{
				ClientID:     clientId,
				ClientSecret: secret,
				TokenURL:     client.TokenEndpoint + "/oauth/token",
			}

			client.httpClient.Transport = &oauth2.Transport{
				Source: conf.TokenSource(client.ctx),
				Base:   client.httpClient.Transport,
			}

			return nil
		},
	}
}

func PasswordCredentials(clientId, username, password string) option {
	return option{
		apply: func(client *Client) error {
			conf := &oauth2.Config{
				ClientID:     clientId,
				ClientSecret: "",
				Endpoint: oauth2.Endpoint{
					AuthURL:  client.AuthEndpoint + "/oauth/auth",
					TokenURL: client.TokenEndpoint + "/oauth/token",
				},
			}

			ctx := context.WithValue(client.ctx, oauth2.HTTPClient, client.httpClient)

			token, err := conf.PasswordCredentialsToken(ctx, username, password)

			if err != nil {
				return err
			}

			client.httpClient.Transport = &oauth2.Transport{
				Source: conf.TokenSource(client.ctx, token),
				Base:   client.httpClient.Transport,
			}

			return nil
		},
	}
}
