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
	shellRequest   string
	shellResponse  string
	githubRequest  string
	githubResponse string
	prRequest      string
	prResponse     string
	userRequest    string
	userResponse   string
}

func (f *fakeCalidumClient) SendDiscordRpcRequest(context.Context, []byte) error {
	return nil
}

func (f *fakeCalidumClient) SendEmailRpcRequest(context.Context, []byte) error {
	return nil
}

func (f *fakeCalidumClient) SendShellRpcRequest(_ context.Context, body []byte) ([]byte, error) {
	f.shellRequest = string(body)
	return []byte(f.shellResponse), nil
}

func (f *fakeCalidumClient) SendGithubRpcRequest(_ context.Context, body []byte) ([]byte, error) {
	f.githubRequest = string(body)
	return []byte(f.githubResponse), nil
}

func (f *fakeCalidumClient) SendGithubPRRpcRequest(_ context.Context, body []byte) ([]byte, error) {
	f.prRequest = string(body)
	return []byte(f.prResponse), nil
}

func (f *fakeCalidumClient) SendCedilleUserRpcRequest(_ context.Context, body []byte) ([]byte, error) {
	f.userRequest = string(body)
	return []byte(f.userResponse), nil
}

func (f *fakeCalidumClient) SendClusterRpcRequest(context.Context, []byte) ([]byte, error) {
	return nil, nil
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

func TestGithubPostRequestRoutesToDeploymentRPC(t *testing.T) {
	const apiKey = "test-api-key"
	t.Setenv(ENV_CALIDUM_ROTAE_SERVICE_API_KEY, apiKey)

	client := &fakeCalidumClient{githubResponse: "https://github.com/clubCedille/k8s-shared/actions/runs/1"}
	v := viper.New()
	v.Set(config.FlagAllowedDomains, []string{"https://cedille.etsmtl.ca"})
	router := initHTTPServerHandler(context.Background(), v, client)

	requestBody := `{"workflow":"grav","clubName":"calidum-club","userUID":"2","domain":"calidum.etsmtl.ca"}`
	request := httptest.NewRequest(http.MethodPost, GITHUB_POST_REQUEST, strings.NewReader(requestBody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-API-KEY", apiKey)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("POST /github status = %d, want %d: %s", response.Code, http.StatusOK, response.Body.String())
	}

	var body struct {
		Response string `json:"response"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("POST /github returned invalid JSON: %v", err)
	}
	if body.Response != client.githubResponse {
		t.Fatalf("POST /github response = %q, want %q", body.Response, client.githubResponse)
	}
}

func TestGithubPRPostRequestRoutesToFetchPRRPC(t *testing.T) {
	const apiKey = "test-api-key"
	t.Setenv(ENV_CALIDUM_ROTAE_SERVICE_API_KEY, apiKey)

	client := &fakeCalidumClient{prResponse: `[{"number":12,"title":"PR for user 42"}]`}
	v := viper.New()
	v.Set(config.FlagAllowedDomains, []string{"https://cedille.etsmtl.ca"})
	router := initHTTPServerHandler(context.Background(), v, client)

	requestBody := `{"userUID":"42"}`
	request := httptest.NewRequest(http.MethodPost, GITHUB_PR_POST_REQUEST, strings.NewReader(requestBody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-API-KEY", apiKey)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("POST %s status = %d, want %d: %s", GITHUB_PR_POST_REQUEST, response.Code, http.StatusOK, response.Body.String())
	}

	var body struct {
		Response string `json:"response"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("POST %s returned invalid JSON: %v", GITHUB_PR_POST_REQUEST, err)
	}
	if body.Response != client.prResponse {
		t.Fatalf("POST %s response = %q, want %q", GITHUB_PR_POST_REQUEST, body.Response, client.prResponse)
	}
	if client.prRequest != requestBody {
		t.Fatalf("POST %s forwarded body = %q, want %q", GITHUB_PR_POST_REQUEST, client.prRequest, requestBody)
	}
}

func TestCedilleUserPostRequestReturnsResponse(t *testing.T) {
	const apiKey = "test-api-key"
	t.Setenv(ENV_CALIDUM_ROTAE_SERVICE_API_KEY, apiKey)

	client := &fakeCalidumClient{userResponse: "https://github.com/clubCedille/k8s-shared/actions/runs/2"}
	v := viper.New()
	v.Set(config.FlagAllowedDomains, []string{"https://cedille.etsmtl.ca"})
	router := initHTTPServerHandler(context.Background(), v, client)

	requestBody := `{"githubUsername":"alexvegas22","githubEmail":"alex@cedille.club","teamSre":false,"clusterRole":"Reader","netdataRole":"observer","userUID":"2"}`
	request := httptest.NewRequest(http.MethodPost, GITHUB_USER_POST_REQUEST, strings.NewReader(requestBody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-API-KEY", apiKey)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("POST /github/user status = %d, want %d: %s", response.Code, http.StatusOK, response.Body.String())
	}

	var body struct {
		Response string `json:"response"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("POST /github/user returned invalid JSON: %v", err)
	}
	if body.Response != client.userResponse {
		t.Fatalf("POST /github/user response = %q, want %q", body.Response, client.userResponse)
	}
	if client.userRequest != requestBody {
		t.Fatalf("POST /github/user forwarded body = %q, want %q", client.userRequest, requestBody)
	}
}
