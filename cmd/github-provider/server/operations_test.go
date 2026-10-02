package server

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	github_provider "github.com/clubcedille/calidum-rotae-backend/pkg/proto-gen/github-provider"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestRequestDeployment(t *testing.T) {
	t.Setenv(ENV_GITHUB_TOKEN, "test-token")

	originalClient := http.DefaultClient
	http.DefaultClient = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				t.Fatalf("expected POST request, got %s", req.Method)
			}
			if req.URL.String() != OUTLINE_WORKFLOW_URL {
				t.Fatalf("expected URL %s, got %s", OUTLINE_WORKFLOW_URL, req.URL)
			}
			if got := req.Header.Get("Authorization"); got != "Bearer test-token" {
				t.Fatalf("expected bearer token, got %q", got)
			}

			var payload workflowDispatchBody
			if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
				t.Fatalf("decode dispatch payload: %v", err)
			}
			if payload.Ref != defaultWorkflowRef {
				t.Fatalf("expected ref %q, got %q", defaultWorkflowRef, payload.Ref)
			}
			if !payload.ReturnRunDetails {
				t.Fatal("expected return_run_details to be true")
			}
			if payload.Inputs["nom_club"] != "test-club" {
				t.Fatalf("unexpected nom_club input: %q", payload.Inputs["nom_club"])
			}
			if payload.Inputs["domaine"] != "wiki.test-club.etsmtl.club" {
				t.Fatalf("unexpected domaine input: %q", payload.Inputs["domaine"])
			}

			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body: io.NopCloser(bytes.NewBufferString(
					`{"workflow_run_id":29461610387,"run_url":"https://api.github.com/repos/ClubCedille/k8s-shared/actions/runs/29461610387","html_url":"https://github.com/ClubCedille/k8s-shared/actions/runs/29461610387"}`,
				)),
				Request: req,
			}, nil
		}),
	}
	t.Cleanup(func() {
		http.DefaultClient = originalClient
	})

	response, err := NewServer().RequestDeployment(context.Background(), &github_provider.DeploymentRequest{
		ClubName: "test-club",
		Domain:   "wiki.test-club.etsmtl.club",
		Workflow: "outline",
		UserUID:  "local-test",
	})
	if err != nil {
		t.Fatalf("RequestDeployment returned an error: %v", err)
	}
	if response.GetWorkflowRunID() != 29461610387 {
		t.Fatalf("expected workflow run ID 29461610387, got %d", response.GetWorkflowRunID())
	}
	if response.GetWorkflowRunUrl() != "https://api.github.com/repos/ClubCedille/k8s-shared/actions/runs/29461610387" {
		t.Fatalf("unexpected workflow run URL: %q", response.GetWorkflowRunUrl())
	}
}
