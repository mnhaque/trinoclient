package client

import (
	"github.com/mnhaque/trinoclient/trino-sdk/config"
	"github.com/mnhaque/trinoclient/trino-sdk/request"
)

type Client struct {
	Config *config.Configuration
}

func New(url string, clientName string) *Client {
	client := &Client{
		Config: config.New(url, clientName),
	}
	return client
}

func (c *Client) NewRequest(query string) *request.Response {
	return request.New(c.Config, query)
}
