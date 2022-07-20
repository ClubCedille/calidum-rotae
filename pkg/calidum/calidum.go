package calidum

import (
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
	// & other microservices
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
