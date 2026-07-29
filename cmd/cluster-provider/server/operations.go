package server

import (
	"context"
	"fmt"
	cluster_provider "github.com/clubcedille/calidum-rotae-backend/pkg/proto-gen/cluster-provider"
)

type Server struct {
	cluster_provider.ClusterProviderServer
}

func NewServer() *Server {
	return &Server{}
}

func (server *Server) GetUserResources(ctx context.Context, message *cluster_provider.UserResourceRequest) (*cluster_provider.GetUserResourcesResponse, error) {
	user := message.UserUID
	resp := "TODO"
	fmt.Println("Message", user)

	return &cluster_provider.GetUserResourcesResponse{ UserResourcesResponse: resp}, nil
}
