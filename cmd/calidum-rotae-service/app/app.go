package app

import (
	"context"
	"fmt"

	"github.com/clubcedille/calidum-rotae-backend/cmd/calidum-rotae-service/client"
	"github.com/clubcedille/calidum-rotae-backend/cmd/calidum-rotae-service/server"
	"github.com/clubcedille/calidum-rotae-backend/pkg/calidum"
	discord_provider "github.com/clubcedille/calidum-rotae-backend/pkg/proto-gen/discord-provider"
	email_provider "github.com/clubcedille/calidum-rotae-backend/pkg/proto-gen/email-provider"
    shell_provider "github.com/clubcedille/calidum-rotae-backend/pkg/proto-gen/shell-provider"
	serverutils "github.com/clubcedille/server-utils"
	"github.com/spf13/viper"
)

type CalidumRotaeService struct {
	// Calidum rotae service client
	calidumService calidum.CalidumClient

	// Discord provider service client
	discordProvider discord_provider.DiscordProviderClient

	// Email provider service client
	emailProvider email_provider.EmailProviderClient

    // Shell provider service client
    shellProvider shell_provider.ShellProviderClient

	// HTTP server - our REST API
	httpServer serverutils.Server
}

func InitFromViper(ctx context.Context, v *viper.Viper) (service *CalidumRotaeService, err error) {
	// Create new instance of the service
	service = &CalidumRotaeService{}

	// Create clients of the gRPC providers
	service.discordProvider, service.emailProvider, err = client.InitFromViper(ctx, v)
	if err != nil {
		return nil, fmt.Errorf("error when initializing client providers: %s", err)
	}

	// Initialize calidum rotae service
	service.calidumService, err = service.initService(ctx, v)
	if err != nil {
		return nil, fmt.Errorf("error initializing calidum rotae service: %s", err)
	}

	// Create instance of the HTTP server
	service.httpServer, err = server.InitHTTPServerFromViper(ctx, v, service.calidumService)
	if err != nil {
		return nil, fmt.Errorf("error when initializing the HTTP server: %s", err)
	}

	return
}

func (c *CalidumRotaeService) initService(ctx context.Context, v *viper.Viper) (calidumService *calidum.CalidumService, err error) {
	// Build new calidum rotae service with its dependencies
	calidumService = calidum.NewCalidumService(calidum.Dependencies{
		DiscordProviderService: c.discordProvider,
		EmailProviderService:   c.emailProvider,
	})

	return
}

func (c *CalidumRotaeService) Run(ctx context.Context, port int32) error {
	return c.httpServer.Run(ctx, serverutils.RunRequest{
		Port:              port,
		ShutdownTimeoutMs: 10000, // 10 seconds
	})
}
