package server

import (
	"context"
	"fmt"
	github_provider "github.com/clubcedille/calidum-rotae-backend/pkg/proto-gen/github-provider"
)

type Server struct {
	github_provider.GithubProviderServer
}

func NewServer() *Server {
	return &Server{}
}

func (server *Server) FetchPR(ctx context.Context, message *github_provider.FetchPRRequest) (*github_provider.FetchPRResponse, error) {
	user := message.UserUID
	commandResponse := "todo"
	fmt.Println("Message", user)

	return &github_provider.FetchPRResponse{ ActiveUserRequests: commandResponse}, nil
}
