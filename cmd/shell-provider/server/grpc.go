package server

import (
	shell_provider "github.com/clubcedille/calidum-rotae-backend/pkg/proto-gen/shell-provider"
	"google.golang.org/grpc"
)

func (s *Server) ConfigureGrpc() *grpc.Server {
	// Create new gRPC server
	grpcServer := grpc.NewServer()

	// Register gRPC providers to the server
	shell_provider.RegisterShellProviderServer(grpcServer, s)

	// Return created server instance
	return grpcServer
}
