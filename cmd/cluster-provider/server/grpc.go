package server

import (
	cluster_provider "github.com/clubcedille/calidum-rotae-backend/pkg/proto-gen/cluster-provider"
	"google.golang.org/grpc"
)

func (s *Server) ConfigureGrpc() *grpc.Server {
	// Create new gRPC server
	grpcServer := grpc.NewServer()

	// Register gRPC providers to the server
	cluster_provider.RegisterClusterProviderServer(grpcServer, s)

	// Return created server instance
	return grpcServer
}
