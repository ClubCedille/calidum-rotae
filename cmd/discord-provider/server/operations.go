package server

import (
	"context"
	"log"

	discord_provider "github.com/clubcedille/calidum-rotae-backend/pkg/proto-gen/discord-provider"
)

type Server struct {
	discord_provider.DiscordProviderServer
}

func NewServer() *Server {
	return &Server{}
}

func (server *Server) SendMessage(ctx context.Context, message *discord_provider.SendMessageRequest) (*discord_provider.SendMessageResponse, error) {
	// Sends a message to a Discord webhook here

	log.Printf("Received message content from client: %s", message.Sender)
	return &discord_provider.SendMessageResponse{}, nil
}
