package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	context "context"

	email_provider "github.com/clubcedille/calidum-rotae-backend/pkg/proto-gen/email-provider"
)

type Server struct {
	email_provider.EmailProviderServer
}

func (server *Server) SendMessage(ctx context.Context, message *email_provider.SendEmailRequest) (*email_provider.SendEmailResponse, error) {
	log.Printf("Received message content from client: %s", message.Sender)
	// Sends a message to a Discord webhook here
	return &email_provider.SendEmailResponse{}, nil
}

const PORT string = "PORT"

func main() {
	lt, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("Error starting tcp listener on %s : %v", PORT, err)
	}

	server := Server{}
	grpcServer := grpc.NewServer()
	email_provider.RegisterEmailProviderServer(grpcServer, &server)

	if err := grpcServer.Serve(lt); err != nil {
		log.Fatalf("Failed to serve gRPC server over port %s : %v", PORT, err)
	}
}
