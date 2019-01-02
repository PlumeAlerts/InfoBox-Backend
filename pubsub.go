package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PlumeAlerts/InfoBoxes-Backend/jwt"
	jwtgo "github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"net/http"
	"time"
)

const url = "https://api.twitch.tv/extensions/message/"

type PubSub struct {
	ContentType string   `json:"content_type"`
	Message     string   `json:"message"`
	Targets     []string `json:"targets"`
}

type Message struct {
	Title       string `json:"title"`
	Icon        string `json:"icon"`
	IconURL     string `json:"iconURL"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Duration    int    `json:"duration"`
}

func SendPubSubMessage(channelId string, msg Message) {
	claims := jwt.JWTClaims{
		ChannelID: channelId,
		Role:      "external",
		Permissions: jwt.JWTPubSubPermissions{
			Send: []string{"broadcast"},
		},
		StandardClaims: jwtgo.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute*3).UnixNano() / int64(time.Millisecond),
		},
	}
	token := jwt.NewJWT(claims)

	message, _ := json.Marshal(msg)
	var data = PubSub{
		ContentType: "application/json",
		Message:     string(message),
		Targets:     []string{"broadcast"},
	}
	fmt.Println(len(data.Message))
	b, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url+channelId, bytes.NewBuffer(b))
	req.Header.Set("Client-Id", "lnusyooewgjzr3sncje0wf1i6gep3x")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
