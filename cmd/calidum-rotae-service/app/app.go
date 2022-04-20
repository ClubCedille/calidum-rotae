package app

import (
	"context"
	"fmt"

	"github.com/clubcedille/calidum-rotae-backend/cmd/calidum-rotae-service/client"
	"github.com/clubcedille/calidum-rotae-backend/cmd/calidum-rotae-service/server"
	discord_provider "github.com/clubcedille/calidum-rotae-backend/pkg/proto-gen/discord-provider"
	email_provider "github.com/clubcedille/calidum-rotae-backend/pkg/proto-gen/email-provider"
	"github.com/spf13/viper"
)

type CalidumRotaeService struct {
	// Discord provider service client
	discordProvider discord_provider.DiscordProviderClient

	// Email provider service client
	emailProvider email_provider.EmailProviderClient

	// HTTP server - our REST API
	httpServer *server.HTTPServer
}

func InitFromViper(ctx context.Context, v *viper.Viper) (service *CalidumRotaeService, err error) {
	// Create new instance of the service
	service = &CalidumRotaeService{}

	// Create clients of the gRPC providers
	service.discordProvider, service.emailProvider, err = client.InitFromViper(ctx, v)
	if err != nil {
		return nil, fmt.Errorf("error when initializing client providers: %s", err)
	}

	// Create instance of the HTTP server
	service.httpServer, err = server.InitHTTPServerFromViper(ctx, v)
	if err != nil {
		return nil, fmt.Errorf("error when initializing the HTTP server: %s", err)
	}

	return
}

func (c *CalidumRotaeService) Run(ctx context.Context) error {
	return c.httpServer.Serve()
}
