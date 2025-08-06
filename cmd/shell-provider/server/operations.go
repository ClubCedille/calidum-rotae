package server

import (
	"context"
    "strings"
	shell_provider "github.com/clubcedille/calidum-rotae-backend/pkg/proto-gen/shell-provider"
)

type Server struct {
	shell_provider.ShellProviderServer
}

func NewServer() *Server {
	return &Server{}
}

func (server *Server) SendCommand(ctx context.Context, message *shell_provider.SendCommandRequest) (*shell_provider.SendCommandResponse, error) {
    command := message.RequestCommand;
    commandResponse := "Privilèges insuffisants"
    if strings.Contains(command, "sudo") {
        commandResponse = "Défi réussi! réponse : Hello, CEDILLE"
    }
	return &shell_provider.SendCommandResponse{CommandResponse: commandResponse}, nil
}
