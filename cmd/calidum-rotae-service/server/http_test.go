package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/clubcedille/calidum-rotae-backend/cmd/calidum-rotae-service/config"
	"github.com/spf13/viper"
)

type fakeCalidumClient struct {
	shellRequest  string
	shellResponse string
}

func (f *fakeCalidumClient) SendDiscordRpcRequest(context.Context, []byte) error {
	return nil
}

func (f *fakeCalidumClient) SendEmailRpcRequest(context.Context, []byte) error {
	return nil
}

func (f *fakeCalidumClient) SendShellRpcRequest(_ context.Context, body []byte) (string, error) {
	f.shellRequest = string(body)
	return f.shellResponse, nil
}

func TestShellPostRequestReturnsDirectJSONResponse(t *testing.T) {
	const apiKey = "test-api-key"
	t.Setenv(ENV_CALIDUM_ROTAE_SERVICE_API_KEY, apiKey)

	client := &fakeCalidumClient{shellResponse: "Défi réussi!"}
	v := viper.New()
	v.Set(config.FlagAllowedDomains, []string{"https://cedille.etsmtl.ca"})
	router := initHTTPServerHandler(context.Background(), v, client)

	request := httptest.NewRequest(http.MethodPost, "/command", strings.NewReader(`{"requestCommand":"sudo ls"}`))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-API-KEY", apiKey)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("POST /command status = %d, want %d", response.Code, http.StatusOK)
	}

	var body struct {
		Response string `json:"response"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("POST /command returned invalid JSON: %v", err)
	}
	if body.Response != client.shellResponse {
		t.Fatalf("POST /command response = %q, want %q", body.Response, client.shellResponse)
	}
	if client.shellRequest != `{"requestCommand":"sudo ls"}` {
		t.Fatalf("POST /command forwarded body = %q", client.shellRequest)
	}
}
