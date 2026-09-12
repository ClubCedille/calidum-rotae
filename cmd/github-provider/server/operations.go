package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	github_provider "github.com/clubcedille/calidum-rotae-backend/pkg/proto-gen/github-provider"
)

const (
	ENV_GITHUB_TOKEN = "GITHUB_TOKEN"
	OUTLINE_WORKFLOW_URL = "https://api.github.com/repos/clubCedille/k8s-shared/actions/workflows/request-outline.yml/dispatches"
	GRAV_WORKFLOW_URL = "https://api.github.com/repos/clubCedille/k8s-shared/actions/workflows/request-grav.yml/dispatches"
)

type RequestOutlineBody struct {
	Ref    string        `json:"ref,omitempty"`
	Inputs OutlineInputs `json:"inputs,omitempty"`
}

type OutlineInputs struct {
	NomClub string `json:"nom_club,omitempty"`
}

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

func (server *Server) RequestOutline(ctx context.Context, message *github_provider.OutlineRequest) (*github_provider.OutlineResponse, error) {
	token, found := os.LookupEnv(ENV_GITHUB_TOKEN)

	if !found {
		return &github_provider.OutlineResponse{}, fmt.Errorf("Error getting env var %s", ENV_GITHUB_TOKEN)
	}

	client := &http.Client {
		CheckRedirect: redirectPolicyFunc,
	}

	outlineRequest := RequestOutlineBody{
		Ref: "main", 
		Inputs: OutlineInputs{
			NomClub: message.GetClubName(),
		}
	}
	
	payload := new(bytes.Buffer)
	
	err := json.NewEncoder(payload).Encode(outlineRequest)

	if err != nil {
		return &github_provider.OutlineResponse{}, fmt.Errorf("Error encoding the message %v", outlineRequest)
	}

	req,err := http.NewRequest("POST", OUTLINE_WORKFLOW_URL, payload)

	if err != nil {
		return &github_provider.OutlineResponse{}, fmt.Errorf("Error forming request")
	}

	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("X-Github-Api-Version", "2026-03-10")

	resp, err := client.Do(req)

	if err != nil {
		return &github_provider.OutlineResponse{}, fmt.Errorf("HTTP Request Failed : %v", err)
	}

	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return &github_provider.OutlineResponse{}, fmt.Errorf("Error calling workflow, status code %d\n body: %s", resp.StatusCode, string(respBody))
	}

	log.Printf("Outline Workflow triggered")
	return &github_provider.OutlineResponse{}, nil


}
