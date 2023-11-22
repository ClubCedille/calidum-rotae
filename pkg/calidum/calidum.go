package calidum

import (
	"context"
	"encoding/json"
	"fmt"

	instrumentation "github.com/clubcedille/calidum-rotae-backend/cmd/otel-instrumentation"
	discord_provider "github.com/clubcedille/calidum-rotae-backend/pkg/proto-gen/discord-provider"
	email_provider "github.com/clubcedille/calidum-rotae-backend/pkg/proto-gen/email-provider"
)

type CalidumService struct {
	discordProviderService discord_provider.DiscordProviderClient
	emailProviderService   email_provider.EmailProviderClient
	tracers                instrumentation.Tracer
}

type CalidumClient interface {
	SendDiscordRpcRequest(ctx context.Context, body []byte) (err error)
	SendEmailRpcRequest(ctx context.Context, body []byte) (err error)
}

var _ CalidumClient = &CalidumService{}

type Dependencies struct {
	DiscordProviderService discord_provider.DiscordProviderClient
	EmailProviderService   email_provider.EmailProviderClient
	Tracer                 instrumentation.Tracer
}

func NewCalidumService(deps Dependencies) *CalidumService {
	return &CalidumService{
		discordProviderService: deps.DiscordProviderService,
		emailProviderService:   deps.EmailProviderService,
		tracers:                deps.Tracer,
	}
}

func (c *CalidumService) SendDiscordRpcRequest(ctx context.Context, body []byte) (err error) {
	_, span := c.tracers.StartSpanAndInitTracers(ctx, instrumentation.DISCORD_PROVIDER_SERVICE)
	defer span.End()

	var data *discord_provider.SendMessageRequest
	if err = json.Unmarshal(body, &data); err != nil {
		errorMessage := fmt.Sprintf("error binding JSON data to gRPC discord object: %s", err.Error())
		span = c.tracers.AddErrorEvent(ctx, span, errorMessage)
		return fmt.Errorf(errorMessage)
	}

	resp, err := c.discordProviderService.SendMessage(ctx, &discord_provider.SendMessageRequest{
		Sender:         data.Sender,
		RequestService: data.RequestService,
		RequestDetails: data.RequestDetails,
	})
	if err != nil {
		errorMessage := fmt.Sprintf("error sending rpc request to discord provider: %s Response: %s", err.Error(), resp)
		span = c.tracers.AddErrorEvent(ctx, span, errorMessage)
		return fmt.Errorf(errorMessage)
	}

	c.tracers.EndSpan(span, "discord message sent")
	return nil
}

func (c *CalidumService) SendEmailRpcRequest(ctx context.Context, body []byte) (err error) {
	_, span := c.tracers.StartSpanAndInitTracers(ctx, instrumentation.EMAIL_PROVIDER_SERVICE)
	defer span.End()

	var data *email_provider.SendEmailRequest
	if err = json.Unmarshal(body, &data); err != nil {
		errorMessage := fmt.Sprintf("error binding JSON data to gRPC email object: %s", err.Error())
		span = c.tracers.AddErrorEvent(ctx, span, errorMessage)
		return fmt.Errorf("error binding JSON data to gRPC email object: %s", err.Error())
	}

	resp, err := c.emailProviderService.SendEmail(ctx, &email_provider.SendEmailRequest{
		Sender:         data.Sender,
		RequestService: data.RequestService,
		RequestDetails: data.RequestDetails,
	})
	if err != nil {
		errorMessage := fmt.Sprintf("error sending rpc request to discord provider: %s Response: %s", err.Error(), resp)
		span = c.tracers.AddErrorEvent(ctx, span, errorMessage)
		return fmt.Errorf(errorMessage)
	}

	c.tracers.EndSpan(span, "email sent")
	return nil
}
