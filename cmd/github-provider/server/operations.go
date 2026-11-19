package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	github_provider "github.com/clubcedille/calidum-rotae-backend/pkg/proto-gen/github-provider"
)

const (
	ENV_GITHUB_TOKEN          = "GITHUB_TOKEN"
	GITHUB_API_URL            = "https://api.github.com"
	K8S_SHARED_REPOSITORY     = "ClubCedille/k8s-shared"
	OUTLINE_WORKFLOW_URL      = "https://api.github.com/repos/clubCedille/k8s-shared/actions/workflows/request-outline.yml/dispatches"
	GRAV_WORKFLOW_URL         = "https://api.github.com/repos/clubCedille/k8s-shared/actions/workflows/request-grav.yml/dispatches"
	S3_WORKFLOW_URL           = "https://api.github.com/repos/clubCedille/k8s-shared/actions/workflows/request-s3.yml/dispatches"
	CEDILLE_USER_WORKFLOW_URL = "https://api.github.com/repos/clubCedille/Plateforme-Cedille/actions/workflows/add-new-member.yml/dispatches"

	pullRequestsPerPage = 100
)

const defaultWorkflowRef = "main"

type workflowDispatchBody struct {
	Ref              string            `json:"ref,omitempty"`
	Inputs           map[string]string `json:"inputs,omitempty"`
	ReturnRunDetails bool              `json:"return_run_details"`
}

type workflowDispatchResponse struct {
	WorkflowRunId int64  `json:"workflow_run_id,omitempty"`
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
	pullRequests, err := fetchPullRequestsByUser(ctx, message.GetUserUID())
	if err != nil {
		return &github_provider.FetchPRResponse{}, err
	}

	response, err := json.Marshal(pullRequests)
	if err != nil {
		return &github_provider.FetchPRResponse{}, fmt.Errorf("error encoding pull requests: %s", err)
	}

	return &github_provider.FetchPRResponse{ActiveUserRequests: string(response)}, nil
}

func (server *Server) RequestDeployment(ctx context.Context, message *github_provider.DeploymentRequest) (*github_provider.WorkflowResponse, error) {
	workflowURLs := map[string]string{
		"outline": OUTLINE_WORKFLOW_URL,
		"grav":    GRAV_WORKFLOW_URL,
		"s3":      S3_WORKFLOW_URL,
	}

	url, found := workflowURLs[message.GetWorkflow()]
	if !found {
		return &github_provider.WorkflowResponse{}, fmt.Errorf("unknown workflow: %s", message.GetWorkflow())
	}

	resp, err := dispatchWorkflow(ctx, url, deploymentWorkflowInputs(message))
	if err != nil {
		return &github_provider.WorkflowResponse{}, err
	}

	log.Printf("%s workflow triggered", message.GetWorkflow())

	var workflowResp workflowDispatchResponse
	err = json.Unmarshal([]byte(resp), &workflowResp)

	if err != nil {
		return &github_provider.WorkflowResponse{}, err
	}
	return &github_provider.WorkflowResponse{WorkflowRunUrl: workflowResp.RunUrl, WorkflowRunID: workflowResp.WorkflowRunId}, nil
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
	return &github_provider.WorkflowResponse{WorkflowRunUrl: workflowResp.RunUrl, WorkflowRunID: workflowResp.WorkflowRunId}, nil
}

func fetchPullRequestsByUser(ctx context.Context, userUID string) ([]json.RawMessage, error) {
	if userUID == "" {
		return nil, fmt.Errorf("user UID is required")
	}

	pullRequests := make([]json.RawMessage, 0)
	for page := 1; ; page++ {
		query := url.Values{}
		query.Set("state", "all")
		query.Set("per_page", strconv.Itoa(pullRequestsPerPage))
		query.Set("page", strconv.Itoa(page))
		query.Set("labels", "userid:"+userUID)

		requestURL := fmt.Sprintf("%s/repos/%s/issues?%s", GITHUB_API_URL, K8S_SHARED_REPOSITORY, query.Encode())
		req, err := newGithubRequest(ctx, http.MethodGet, requestURL, nil)
		if err != nil {
			return nil, err
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("github pull request request failed: %s", err)
		}

		body, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			return nil, fmt.Errorf("error reading github pull request response: %s", readErr)
		}
		if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
			return nil, fmt.Errorf("github pull request request failed with status %d: %s", resp.StatusCode, string(body))
		}

		var issues []json.RawMessage
		if err := json.Unmarshal(body, &issues); err != nil {
			return nil, fmt.Errorf("error decoding github pull request response: %s", err)
		}

		for _, issue := range issues {
			var metadata struct {
				PullRequest *json.RawMessage `json:"pull_request"`
			}
			if err := json.Unmarshal(issue, &metadata); err != nil {
				return nil, fmt.Errorf("error decoding github issue response: %s", err)
			}
			if metadata.PullRequest != nil {
				pullRequests = append(pullRequests, issue)
			}
		}

		if len(issues) < pullRequestsPerPage {
			break
		}
	}

	return pullRequests, nil
}

func deploymentWorkflowInputs(message *github_provider.DeploymentRequest) map[string]string {
	return map[string]string{
		"nom_club": message.GetClubName(),
		"domaine":  message.GetDomain(),
		"id":       message.GetUserUID(),
		"club":     message.GetClubName(),
	}
}

func newGithubRequest(ctx context.Context, method, requestURL string, body io.Reader) (*http.Request, error) {
	token, found := os.LookupEnv(ENV_GITHUB_TOKEN)
	if !found {
		return nil, fmt.Errorf("error getting env var %s", ENV_GITHUB_TOKEN)
	}

	req, err := http.NewRequestWithContext(ctx, method, requestURL, body)
	if err != nil {
		return nil, fmt.Errorf("error forming request: %s", err)
	}
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("X-Github-Api-Version", "2026-03-10")

	return req, nil
}

func dispatchWorkflow(ctx context.Context, url string, inputs map[string]string) (string, error) {
	payload := workflowDispatchBody{
		Ref:              defaultWorkflowRef,
		Inputs:           inputs,
		ReturnRunDetails: true,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(payload); err != nil {
		return "", fmt.Errorf("error encoding the dispatch body %v: %s", payload, err)
	}

	req, err := newGithubRequest(ctx, http.MethodPost, url, &buf)
	if err != nil {
		return "", err
	}

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
