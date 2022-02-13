package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	context "context"

	discord_provider "github.com/clubcedille/calidum-rotae-backend/pkg/proto-gen/discord-provider"
)

type Server struct {
	discord_provider.DiscordProviderServer
}

func (server *Server) SendMessage(ctx context.Context, message *discord_provider.SendMessageRequest) (*discord_provider.SendMessageResponse, error) {
	log.Printf("Received message content from client: %s", message.Sender)
	// Sends a message to a Discord webhook here
	return &discord_provider.SendMessageResponse{}, nil
}

const PORT string = "PORT"

func main() {
	lt, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("Error starting tcp listener on %s : %v", PORT, err)
	}

	server := Server{}
	grpcServer := grpc.NewServer()
	discord_provider.RegisterDiscordProviderServer(grpcServer, &server)

	if err := grpcServer.Serve(lt); err != nil {
		log.Fatalf("Failed to serve gRPC server over port %s : %v", PORT, err)
	}
}
