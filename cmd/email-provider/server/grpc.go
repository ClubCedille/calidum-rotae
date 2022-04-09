package server

import (
	email_provider "github.com/clubcedille/calidum-rotae-backend/pkg/proto-gen/email-provider"
	"google.golang.org/grpc"
)

func (s *Server) ConfigureGrpc() *grpc.Server {
	// Create new gRPC server
	grpcServer := grpc.NewServer()

	// Register gRPC providers to the server
	email_provider.RegisterEmailProviderServer(grpcServer, s)

	// Return created server instance
	return grpcServer
}
