package client

import (
	"context"
	"fmt"

	"github.com/clubcedille/calidum-rotae-backend/cmd/calidum-rotae-service/config"
	discord_provider "github.com/clubcedille/calidum-rotae-backend/pkg/proto-gen/discord-provider"
	email_provider "github.com/clubcedille/calidum-rotae-backend/pkg/proto-gen/email-provider"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func InitFromViper(ctx context.Context, v *viper.Viper) (
	discordClient discord_provider.DiscordProviderClient,
	emailClient email_provider.EmailProviderClient,
	err error,
) {
	discordClient, err = initDiscordProviderClient(ctx, v)
	if err != nil {
		err = fmt.Errorf("error when initializing discord provider client: %s", err)
		return
	}

	emailClient, err = initEmailProviderClient(ctx, v)
	if err != nil {
		err = fmt.Errorf("error when initializing email provider client: %s", err)
		return
	}

	return
}

func initDiscordProviderClient(ctx context.Context, v *viper.Viper) (client discord_provider.DiscordProviderClient, err error) {
	hostname := v.GetString(config.FlagDiscordProviderHostname)
	port := v.GetUint32(config.FlagDiscordProviderPort)
	target := fmt.Sprintf("%s:%d", hostname, port)

	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("error when dialing discord microservice: %s", err)
	}

	return discord_provider.NewDiscordProviderClient(conn), nil
}

func initEmailProviderClient(ctx context.Context, v *viper.Viper) (client email_provider.EmailProviderClient, err error) {
	hostname := v.GetString(config.FlagEmailProviderHostname)
	port := v.GetUint32(config.FlagEmailProviderPort)
	target := fmt.Sprintf("%s:%d", hostname, port)

	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("error when dialing email microservice: %s", err)
	}

	return email_provider.NewEmailProviderClient(conn), nil
}
