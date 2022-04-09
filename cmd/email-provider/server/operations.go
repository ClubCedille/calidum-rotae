package server

import (
	"context"
	"log"

	email_provider "github.com/clubcedille/calidum-rotae-backend/pkg/proto-gen/email-provider"
)

type Server struct {
	email_provider.EmailProviderServer
}

func NewServer() *Server {
	return &Server{}
}

func (server *Server) SendMessage(ctx context.Context, message *email_provider.SendEmailRequest) (*email_provider.SendEmailResponse, error) {
	// TODO: send an email to someone here

	log.Printf("Received message content from client: %s", message.Sender)
	return &email_provider.SendEmailResponse{}, nil
}
