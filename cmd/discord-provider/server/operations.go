package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	discord_provider "github.com/clubcedille/calidum-rotae-backend/pkg/proto-gen/discord-provider"
)

type Server struct {
	discord_provider.DiscordProviderServer
}

func NewServer() *Server {
	return &Server{}
}

func (server *Server) SendMessage(ctx context.Context, message *discord_provider.SendMessageRequest) (*discord_provider.SendMessageResponse, error) {
	url, found := os.LookupEnv("DISCORD_WEBHOOK_URL")
	if !found {
		return &discord_provider.SendMessageResponse{}, fmt.Errorf("the env var DISCORD_WEBHOOK_URL is not set to a value")
	}

	discordName := fmt.Sprintf("**Name:** %s %s", message.Sender.FirstName, message.Sender.LastName)
	discordEmail := fmt.Sprintf("**Email:** %s", message.Sender.Email)
	discordPhone := fmt.Sprintf("**Phone:** %s", message.Sender.PhoneNumber)
	discordRequestDetails := fmt.Sprintf("**Request details:** %s", message.RequestDetails)
	discordRequestServices := fmt.Sprintf("**Request services:** %s", message.RequestService)
	discordEmbedded := Embedded{
		Title:       "New submission",
		Description: fmt.Sprintf("%s\n%s\n%s\n\n%s\n%s", discordName, discordEmail, discordPhone, discordRequestDetails, discordRequestServices),
		Color:       "16745728", // orange
		Footer: Footer{
			Text: "By calidum-rotae services",
		},
	}

	discordMessage := Message{
		Username:  "calidum-rotae",
		Embeddeds: []Embedded{discordEmbedded},
	}

	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(discordMessage)
	resp, err := http.Post(url, "application/json", payload)
	if err != nil {
		return &discord_provider.SendMessageResponse{}, fmt.Errorf(strconv.Itoa(resp.StatusCode))
	}

	return &discord_provider.SendMessageResponse{}, nil
}
