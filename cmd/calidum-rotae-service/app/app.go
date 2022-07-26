package app

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/clubcedille/calidum-rotae-backend/cmd/calidum-rotae-service/client"
	"github.com/clubcedille/calidum-rotae-backend/cmd/calidum-rotae-service/config"
	"github.com/clubcedille/calidum-rotae-backend/cmd/calidum-rotae-service/server"
	"github.com/clubcedille/calidum-rotae-backend/pkg/calidum"
	"github.com/clubcedille/calidum-rotae-backend/pkg/database/postgres"
	discord_provider "github.com/clubcedille/calidum-rotae-backend/pkg/proto-gen/discord-provider"
	email_provider "github.com/clubcedille/calidum-rotae-backend/pkg/proto-gen/email-provider"
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

	// Create instance of the HTTP server
	service.httpServer, err = server.InitHTTPServerFromViper(ctx, v)
	if err != nil {
		return nil, fmt.Errorf("error when initializing the HTTP server: %s", err)
	}

	// Initialize calidum rotae service
	service.calidumService, err = service.initService(ctx, v)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize calidum rotae service: %s", err)
	}

	return
}

func (c *CalidumRotaeService) initService(ctx context.Context, v *viper.Viper) (calidumService *calidum.CalidumService, err error) {
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		return nil, errors.New("failed to fetch the following environment variable: 'DB_PASSWORD' is not set")
	}

	db, err := postgres.NewPostgresClient(postgres.Config{
		User:        v.GetString(config.FlagDbUser),
		Password:    dbPassword,
		DbName:      v.GetString(config.FlagDbName),
		Host:        v.GetString(config.FlagDbHost),
		Schema:      v.GetString(config.FlagDbSchema),
		SSLMode:     v.GetString(config.FlagDbSslMode),
		MaxIdleConn: 100,
		MaxOpenConn: 100,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %s", err)
	}

	// Build new calidum rotae service with its dependencies
	calidumService = calidum.NewCalidumService(calidum.Dependencies{
		Database:               db,
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
