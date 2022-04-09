package server

import (
	discord_provider "github.com/clubcedille/calidum-rotae-backend/pkg/proto-gen/discord-provider"
	"google.golang.org/grpc"
)

func (s *Server) ConfigureGrpc() *grpc.Server {
	// Create new gRPC server
	grpcServer := grpc.NewServer()

	// Register gRPC providers to the server
	discord_provider.RegisterDiscordProviderServer(grpcServer, s)

	// Return created server instance
	return grpcServer
}
