package calidum

import (
	"context"
	"testing"

	shell_provider "github.com/clubcedille/calidum-rotae-backend/pkg/proto-gen/shell-provider"
	"google.golang.org/grpc"
)

type fakeShellProvider struct {
	request  *shell_provider.SendCommandRequest
	response *shell_provider.SendCommandResponse
	err      error
}

func (f *fakeShellProvider) SendCommand(_ context.Context, request *shell_provider.SendCommandRequest, _ ...grpc.CallOption) (*shell_provider.SendCommandResponse, error) {
	f.request = request
	return f.response, f.err
}

func TestSendShellRpcRequest(t *testing.T) {
	provider := &fakeShellProvider{
		response: &shell_provider.SendCommandResponse{CommandResponse: "shell response"},
	}
	service := NewCalidumService(Dependencies{ShellProviderService: provider})

	response, err := service.SendShellRpcRequest(context.Background(), []byte(`{"requestCommand":"sudo ls"}`))
	if err != nil {
		t.Fatalf("SendShellRpcRequest() error = %v", err)
	}
	if response != "shell response" {
		t.Fatalf("SendShellRpcRequest() response = %q, want %q", response, "shell response")
	}
	if provider.request.GetRequestCommand() != "sudo ls" {
		t.Fatalf("SendShellRpcRequest() command = %q, want %q", provider.request.GetRequestCommand(), "sudo ls")
	}
}

func TestSendShellRpcRequestRejectsMalformedJSON(t *testing.T) {
	service := NewCalidumService(Dependencies{ShellProviderService: &fakeShellProvider{}})

	if _, err := service.SendShellRpcRequest(context.Background(), []byte(`{"requestCommand":`)); err == nil {
		t.Fatal("SendShellRpcRequest() error = nil, want malformed JSON error")
	}
}
