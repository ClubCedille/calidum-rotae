package server

import (
	"context"
	"testing"

	shell_provider "github.com/clubcedille/calidum-rotae-backend/pkg/proto-gen/shell-provider"
)

func TestSendCommand(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		expected string
	}{
		{name: "empty command", command: "   ", expected: emptyCommandResponse},
		{name: "regular command", command: "ls", expected: insufficientAccessResponse},
		{name: "sudo challenge", command: "sudo ls", expected: challengeSuccessfulResponse},
	}

	server := NewServer()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := server.SendCommand(context.Background(), &shell_provider.SendCommandRequest{
				RequestCommand: tt.command,
			})
			if err != nil {
				t.Fatalf("SendCommand() error = %v", err)
			}
			if response.GetCommandResponse() != tt.expected {
				t.Fatalf("SendCommand() response = %q, want %q", response.GetCommandResponse(), tt.expected)
			}
		})
	}
}
