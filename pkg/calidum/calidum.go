package calidum

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/clubcedille/calidum-rotae-backend/pkg/database"
	discord_provider "github.com/clubcedille/calidum-rotae-backend/pkg/proto-gen/discord-provider"
	email_provider "github.com/clubcedille/calidum-rotae-backend/pkg/proto-gen/email-provider"
)

type CalidumService struct {
	db                     database.Operations
	discordProviderService discord_provider.DiscordProviderClient
	emailProviderService   email_provider.EmailProviderClient
}

type CalidumClient interface {
	// TODO: Define methods to speak with database
	SendDiscordRpcRequest(ctx context.Context, body []byte) (err error)
	SendEmailRpcRequest(ctx context.Context, body []byte) (err error)
}

var _ CalidumClient = &CalidumService{}

type Dependencies struct {
	Database               database.Operations
	DiscordProviderService discord_provider.DiscordProviderClient
	EmailProviderService   email_provider.EmailProviderClient
}

func NewCalidumService(deps Dependencies) *CalidumService {
	return &CalidumService{
		db:                     deps.Database,
		discordProviderService: deps.DiscordProviderService,
		emailProviderService:   deps.EmailProviderService,
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
