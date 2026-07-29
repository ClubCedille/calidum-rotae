package calidum

import (
	"context"
	"encoding/json"
	"fmt"

	discord_provider "github.com/clubcedille/calidum-rotae-backend/pkg/proto-gen/discord-provider"
	email_provider "github.com/clubcedille/calidum-rotae-backend/pkg/proto-gen/email-provider"
	shell_provider "github.com/clubcedille/calidum-rotae-backend/pkg/proto-gen/shell-provider"
	github_provider "github.com/clubcedille/calidum-rotae-backend/pkg/proto-gen/github-provider"
	cluster_provider "github.com/clubcedille/calidum-rotae-backend/pkg/proto-gen/cluster-provider"
)

type CalidumService struct {
	discordProviderService discord_provider.DiscordProviderClient
	emailProviderService   email_provider.EmailProviderClient
	shellProviderService   shell_provider.ShellProviderClient
	githubProviderService   github_provider.GithubProviderClient
	clusterProviderService   cluster_provider.ClusterProviderClient
}

type CalidumClient interface {
	SendDiscordRpcRequest(ctx context.Context, body []byte) (err error)
	SendEmailRpcRequest(ctx context.Context, body []byte) (err error)
	SendShellRpcRequest(ctx context.Context, body []byte) (response []byte, err error)
	SendGithubRpcRequest(ctx context.Context, body []byte) (response []byte, err error)
	SendClusterRpcRequest(ctx context.Context, body []byte) (response []byte, err error)
}

var _ CalidumClient = &CalidumService{}

type Dependencies struct {
	DiscordProviderService discord_provider.DiscordProviderClient
	EmailProviderService   email_provider.EmailProviderClient
	ShellProviderService   shell_provider.ShellProviderClient
	GithubProviderService   github_provider.GithubProviderClient
	ClusterProviderService   cluster_provider.ClusterProviderClient
}

func NewCalidumService(deps Dependencies) *CalidumService {
	return &CalidumService{
		discordProviderService: deps.DiscordProviderService,
		emailProviderService:   deps.EmailProviderService,
		shellProviderService:   deps.ShellProviderService,
		githubProviderService:   deps.GithubProviderService,
		clusterProviderService:   deps.ClusterProviderService,
	}
}

func (c *CalidumService) SendDiscordRpcRequest(ctx context.Context, body []byte) (err error) {
	var data *discord_provider.SendMessageRequest
	if err = json.Unmarshal(body, &data); err != nil {
		return fmt.Errorf("error binding JSON data to gRPC discord object: %s", err.Error())
	}

	resp, err := c.discordProviderService.SendMessage(ctx, &discord_provider.SendMessageRequest{
		Sender:         data.Sender,
		RequestService: data.RequestService,
		RequestDetails: data.RequestDetails,
	})
	if err != nil {
		return fmt.Errorf("error sending rpc request to discord provider: %s Response: %s", err.Error(), resp)
	}

	return nil
}

func (c *CalidumService) SendEmailRpcRequest(ctx context.Context, body []byte) (err error) {
	var data *email_provider.SendEmailRequest
	if err = json.Unmarshal(body, &data); err != nil {
		return fmt.Errorf("error binding JSON data to gRPC email object: %s", err.Error())
	}

	resp, err := c.emailProviderService.SendEmail(ctx, &email_provider.SendEmailRequest{
		Sender:         data.Sender,
		RequestService: data.RequestService,
		RequestDetails: data.RequestDetails,
	})
	if err != nil {
		return fmt.Errorf("error sending rpc request to email provider: %s Response: %s", err.Error(), resp)
	}

	return nil
}

func (c *CalidumService) SendShellRpcRequest(ctx context.Context, body []byte) (jsonResponse []byte, err error) {
	var data *shell_provider.SendCommandRequest
	if err = json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("error binding JSON data to gRPC shell object: %s", err.Error())
	}

	resp, err := c.shellProviderService.SendCommand(ctx, &shell_provider.SendCommandRequest{
		RequestCommand: data.RequestCommand,
	})
	if err != nil {
		return nil, fmt.Errorf("error sending rpc request to shell provider: %s", err.Error())
	}

	// Convert the response to JSON
	jsonResponse, err = json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("error converting response to JSON: %s", err.Error())
	}

	return jsonResponse, nil
}

func (c *CalidumService) SendGithubRpcRequest(ctx context.Context, body []byte) (jsonResponse []byte, err error) {
	var data *github_provider.FetchPRRequest
	jsonResponse = []byte("")
	if err = json.Unmarshal(body, &data); err != nil {
		return jsonResponse, fmt.Errorf("error binding JSON data to gRPC Github object: %s", err.Error())
	}

	resp, err := c.githubProviderService.FetchPR(ctx, &github_provider.FetchPRRequest{
		UserUID:         data.UserUID,
	})
	jsonResponse = []byte(resp.String())
	if err != nil {
		return jsonResponse, fmt.Errorf("error sending rpc request to Github provider: %s Response: %s", err.Error(), resp)
	}

	return jsonResponse, nil
}

func (c *CalidumService) SendClusterRpcRequest(ctx context.Context, body []byte) (resp []byte, err error) {

	var data *cluster_provider.UserResourceRequest
	if err = json.Unmarshal(body, &data); err != nil {
		return resp, fmt.Errorf("error binding JSON data to gRPC cluster object: %s", err.Error())
	}

	response, err := c.clusterProviderService.GetUserResources(ctx, &cluster_provider.UserResourceRequest{
		UserUID: data.UserUID,
	})
	resp = []byte(response.String())
	if err != nil {
		return resp, fmt.Errorf("error sending rpc request to cluster provider: %s Response: %s", err.Error(), resp)
	}

	return resp, nil
}
