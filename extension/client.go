package extension

import (
	"github.com/dghubble/sling"
)

const (
	ApiBase = "https://api.twitch.tv/extensions/"
)

type Client struct {
	OwnerID      string
	ClientID     string
	ClientSecret string
}

// OwnerID is the extension owner
func NewClient(ownerID string, clientID string, clientSecret string) *Client {
	return &Client{
		OwnerID:      ownerID,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
}

func (c *Client) getAPIBase(hasClientID bool) string {
	if hasClientID {
		return ApiBase + c.ClientID + "/"
	} else {
		return ApiBase
	}
}

func (c *Client) request(jwtToken string) *sling.Sling {
	return c.requestWithClientID(jwtToken, true)
}

func (c *Client) requestWithClientID(jwtToken string, hasClientID bool) *sling.Sling {
	sling := sling.New().Base(c.getAPIBase(hasClientID)).Set("Client-Id", c.ClientID)
	if jwtToken != "" {
		return sling.Set("Authorization", "Bearer "+jwtToken)
	}
	return sling
}
