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
	"strconv"

	github_provider "github.com/clubcedille/calidum-rotae-backend/pkg/proto-gen/github-provider"
)

const (
	ENV_GITHUB_TOKEN          = "GITHUB_TOKEN"
	OUTLINE_WORKFLOW_URL      = "https://api.github.com/repos/clubCedille/k8s-shared/actions/workflows/request-outline.yml/dispatches"
	GRAV_WORKFLOW_URL         = "https://api.github.com/repos/clubCedille/k8s-shared/actions/workflows/request-grav.yml/dispatches"
	CEDILLE_USER_WORKFLOW_URL = "https://api.github.com/repos/clubCedille/Plateforme-Cedille/actions/workflows/add-new-member.yml/dispatches"
)

const defaultWorkflowRef = "main"

type workflowDispatchBody struct {
	Ref    string            `json:"ref,omitempty"`
	Inputs map[string]string `json:"inputs,omitempty"`
}

type workflowDispatchResponse struct {
	WorkflowRunId int `json:"workflow_run_id,omitempty"`
	RunUrl        string `json:"run_url,omitempty"`
	HtmlUrl       string `json:"html_url,omitempty"`
}

type Server struct {
	github_provider.GithubProviderServer
}

func NewServer() *Server {
	return &Server{}
}

func (server *Server) FetchPR(ctx context.Context, message *github_provider.FetchPRRequest) (*github_provider.FetchPRResponse, error) {
	return &github_provider.FetchPRResponse{ActiveUserRequests: "todo"}, nil
}

func (server *Server) RequestDeployment(ctx context.Context, message *github_provider.DeploymentRequest) (*github_provider.WorkflowResponse, error) {
	workflowURLs := map[string]string{
		"outline": OUTLINE_WORKFLOW_URL,
		"grav":    GRAV_WORKFLOW_URL,
	}

	url, found := workflowURLs[message.GetWorkflow()]
	if !found {
		return &github_provider.WorkflowResponse{}, fmt.Errorf("unknown workflow: %s", message.GetWorkflow())
	}

	resp, err := dispatchWorkflow(ctx, url, map[string]string{
		"nom_club": message.GetClubName(),
		"domaine":  message.GetDomain(),
	})
	if err != nil {
		return &github_provider.WorkflowResponse{}, err
	}

	log.Printf("%s workflow triggered", message.GetWorkflow())

	var workflowResp workflowDispatchResponse
    err = json.Unmarshal([]byte(resp), &workflowResp)

	if err != nil {
		return &github_provider.WorkflowResponse{}, err
	}
	return &github_provider.WorkflowResponse{WorkflowRunUrl: workflowResp.RunUrl, WorkflowRunID: int32(workflowResp.WorkflowRunId)}, nil
}

func (server *Server) AddCedilleUser(ctx context.Context, message *github_provider.CedilleUserRequest) (*github_provider.WorkflowResponse, error) {
	resp, err := dispatchWorkflow(ctx, CEDILLE_USER_WORKFLOW_URL, map[string]string{
		"github_username": message.GetGithubUsername(),
		"github_email":    message.GetGithubEmail(),
		"team_sre":        strconv.FormatBool(message.GetTeamSre()),
		"cluster_role":    message.GetClusterRole(),
		"netdata_role":    message.GetNetdataRole(),
	})
	if err != nil {
		return &github_provider.WorkflowResponse{}, err
	}

	log.Printf("Cedille user workflow triggered for %s", message.GetGithubUsername())
	var workflowResp workflowDispatchResponse
    err = json.Unmarshal([]byte(resp), &workflowResp)

	if err != nil {
		return &github_provider.WorkflowResponse{}, err
	}
	return &github_provider.WorkflowResponse{WorkflowRunUrl: workflowResp.RunUrl, WorkflowRunID: int32(workflowResp.WorkflowRunId)}, nil}

func dispatchWorkflow(ctx context.Context, url string, inputs map[string]string) (string, error) {
	token, found := os.LookupEnv(ENV_GITHUB_TOKEN)
	if !found {
		return "", fmt.Errorf("error getting env var %s", ENV_GITHUB_TOKEN)
	}

	payload := workflowDispatchBody{Ref: defaultWorkflowRef, Inputs: inputs}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(payload); err != nil {
		return "", fmt.Errorf("error encoding the dispatch body %v: %s", payload, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &buf)
	if err != nil {
		return "", fmt.Errorf("error forming request: %s", err)
	}
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("X-Github-Api-Version", "2026-03-10")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("http request failed: %s", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading dispatch response: %s", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("workflow dispatch failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return string(respBody), nil
}
