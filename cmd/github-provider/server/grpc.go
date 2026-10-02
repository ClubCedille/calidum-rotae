package server

import (
	github_provider "github.com/clubcedille/calidum-rotae-backend/pkg/proto-gen/github-provider"
	"google.golang.org/grpc"
)

func (s *Server) ConfigureGrpc() *grpc.Server {
	// Create new gRPC server
	grpcServer := grpc.NewServer()

	// Register gRPC providers to the server
	github_provider.RegisterGithubProviderServer(grpcServer, s)

	// Return created server instance
	return grpcServer
}
