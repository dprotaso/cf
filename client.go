package cf

import (
	"context"
	"encoding/json"
	"net/http"
)

type Client struct {
	Endpoint   string
	httpClient *http.Client
	ctx        context.Context

	AuthEndpoint  string
	TokenEndpoint string
}

type optScope int

const (
	afterInfoFetch optScope = iota
	beforeInfoFetch
)

type option struct {
	apply func(*Client) error
	scope optScope
}

func CustomTransport(transport *http.Transport) option {
	return option{
		scope: beforeInfoFetch,
		apply: func(client *Client) error {
			client.httpClient.Transport = transport
			return nil
		},
	}
}

func NewClient(endpoint string, opts ...option) (*Client, error) {
	client := &Client{
		Endpoint:   endpoint,
		httpClient: &http.Client{},
		ctx:        context.Background(),
	}

	if err := applyOptions(client, beforeInfoFetch, opts); err != nil {
		return nil, err
	}

	if err := updateServerInfo(client); err != nil {
		return nil, err
	}

	if err := applyOptions(client, afterInfoFetch, opts); err != nil {
		return nil, err
	}

	return client, nil
}

func applyOptions(c *Client, scope optScope, opts []option) error {
	for _, opt := range opts {
		if opt.scope != scope {
			continue
		}

		if err := opt.apply(c); err != nil {
			return err
		}
	}

	return nil
}

func updateServerInfo(c *Client) error {
	resp, err := c.httpClient.Get(c.Endpoint + "/v2/info")

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	body := struct {
		AuthEndpoint  string `json:"authorization_endpoint"`
		TokenEndpoint string `json:"token_endpoint"`
	}{}

	if err = dec.Decode(&body); err != nil {
		return err
	}

	c.AuthEndpoint = body.AuthEndpoint
	c.TokenEndpoint = body.TokenEndpoint

	return nil
}
