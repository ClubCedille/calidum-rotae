package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	discord_provider "github.com/clubcedille/calidum-rotae-backend/pkg/proto-gen/discord-provider"
)

// environment variables
const (
	ENV_DISCORD_WEBHOOK_URL = "DISCORD_WEBHOOK_URL"
)

type Server struct {
	discord_provider.DiscordProviderServer
}

func NewServer() *Server {
	return &Server{}
}

func (server *Server) SendMessage(ctx context.Context, message *discord_provider.SendMessageRequest) (*discord_provider.SendMessageResponse, error) {
	url, found := os.LookupEnv(ENV_DISCORD_WEBHOOK_URL)
	if !found {
		return &discord_provider.SendMessageResponse{}, fmt.Errorf("error getting env var %s", ENV_DISCORD_WEBHOOK_URL)
	}

	discordMessage := discordMessage{
		Username: "calidum-rotae",
		Embeddeds: []embedded{
			{
				Title: "New submission",
				Description: discordSenderInformation{
					firstName:      message.Sender.FirstName,
					lastName:       message.Sender.LastName,
					email:          message.Sender.Email,
					phoneNumber:    message.Sender.PhoneNumber,
					requestDetails: message.RequestDetails,
					requestService: message.RequestService,
				}.String(),
				Color: "16745728", // orange
				Footer: footer{
					Text: "By calidum-rotae services",
				},
			}},
	}

	payload := new(bytes.Buffer)
	err := json.NewEncoder(payload).Encode(discordMessage)
	if err != nil {
		return &discord_provider.SendMessageResponse{}, fmt.Errorf("error encoding the message %v", discordMessage)
	}

	resp, err := http.Post(url, "application/json", payload)
	if err != nil {
		defer resp.Body.Close()
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return &discord_provider.SendMessageResponse{}, fmt.Errorf("error sending discord webhook: status code %d", resp.StatusCode)
		}
		return &discord_provider.SendMessageResponse{}, fmt.Errorf("error sending discord webhook: status code: %d\n body: %s", resp.StatusCode, string(respBody))

	}

	log.Printf("Discord message sent")
	return &discord_provider.SendMessageResponse{}, nil
}
