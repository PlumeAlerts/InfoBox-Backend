package extension

import (
	"net/http"
)

type PubSubMessageBody struct {
	ContentType string   `json:"content_type"`
	Message     string   `json:"message"`
	Targets     []string `json:"targets"`
}

//POST https://api.twitch.tv/extensions/message/<channel ID>
func (c *Client) PostPubSubMessage(channelID string, message string) error {
	claims := JWTClaims{
		UserID:    c.OwnerID,
		ChannelID: channelID,
		Role:      "external",
		Permissions: JWTPubSubPermissions{
			Send: []string{"broadcast"},
		},
	}
	jwtToken, err := c.NewJWTWithClaim(claims)
	if err != nil {
		return err
	}
	bodyData := PubSubMessageBody{ContentType: "application/json",
		Message: message,
		Targets: []string{"broadcast"}}

	req, err := c.requestWithClientID(jwtToken, false).Post("message/" + channelID).BodyJSON(&bodyData).Request()
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

//POST https://api.twitch.tv/extensions/message/all
func (c *Client) PostPubSubMessageAll(message string) error {
	claims := JWTClaims{
		UserID:    c.OwnerID,
		ChannelID: "all",
		Role:      "external",
		Permissions: JWTPubSubPermissions{
			Send: []string{"global"},
		},
	}
	jwtToken, err := c.NewJWTWithClaim(claims)
	if err != nil {
		return err
	}
	bodyData := PubSubMessageBody{ContentType: "application/json",
		Message: message,
		Targets: []string{"global"}}

	req, err := c.requestWithClientID(jwtToken, false).Post("message/all").BodyJSON(&bodyData).ReceiveSuccess(nil)
	if req.StatusCode == 204 {
		return nil
	}
	return err
}
