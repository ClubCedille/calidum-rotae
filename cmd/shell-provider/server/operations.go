package server

import (
	"context"
	"strings"

	shell_provider "github.com/clubcedille/calidum-rotae-backend/pkg/proto-gen/shell-provider"
)

const (
	emptyCommandResponse        = "Entrez une commande."
	insufficientAccessResponse  = "Privilèges insuffisants"
	challengeSuccessfulResponse = "Défi réussi! réponse : Hello, CEDILLE"
)

type Server struct {
	shell_provider.ShellProviderServer
}

func NewServer() *Server {
	return &Server{}
}

func (server *Server) SendCommand(_ context.Context, message *shell_provider.SendCommandRequest) (*shell_provider.SendCommandResponse, error) {
	command := strings.TrimSpace(message.GetRequestCommand())
	response := insufficientAccessResponse

	if command == "" {
		response = emptyCommandResponse
	} else if strings.Contains(command, "sudo") {
		response = challengeSuccessfulResponse
	}

	return &shell_provider.SendCommandResponse{CommandResponse: response}, nil
}
