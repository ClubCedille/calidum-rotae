package client

import (
	"context"
	"fmt"
	"os"

	"github.com/clubcedille/calidum-rotae-backend/cmd/calidum-rotae-service/config"
	discord_provider "github.com/clubcedille/calidum-rotae-backend/pkg/proto-gen/discord-provider"
	email_provider "github.com/clubcedille/calidum-rotae-backend/pkg/proto-gen/email-provider"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	ENV_CERTIFICATE_FILE_PATH = "CERTIFICATE_FILE_PATH"
)

func InitFromViper(ctx context.Context, v *viper.Viper) (
	discordClient discord_provider.DiscordProviderClient,
	emailClient email_provider.EmailProviderClient,
	err error,
) {

	certFilePath := os.Getenv(ENV_CERTIFICATE_FILE_PATH)

	discordClient, err = initDiscordProviderClient(ctx, v, certFilePath)
	if err != nil {
		err = fmt.Errorf("error when initializing discord provider client: %s", err)
		return
	}

	emailClient, err = initEmailProviderClient(ctx, v, certFilePath)
	if err != nil {
		err = fmt.Errorf("error when initializing email provider client: %s", err)
		return
	}

	return
}

func initDiscordProviderClient(ctx context.Context, v *viper.Viper, certFilePath string) (client discord_provider.DiscordProviderClient, err error) {
	hostname := v.GetString(config.FlagDiscordProviderHostname)
	port := v.GetUint32(config.FlagDiscordProviderPort)
	target := fmt.Sprintf("%s:%d", hostname, port)

	conn := &grpc.ClientConn{}
	if certFilePath != "" {
		creds, err := credentials.NewClientTLSFromFile(certFilePath, "")
		if err != nil {
			return nil, fmt.Errorf("error creating TLS credentials for the discord microservice %v", err)
		}
		conn, err = grpc.Dial(target, grpc.WithTransportCredentials(creds))
		if err != nil {
			return nil, fmt.Errorf("error when dialing discord microservice using TLS: %s", err)
		}
	} else {
		conn, err = grpc.Dial(target, grpc.WithInsecure())
		if err != nil {
			return nil, fmt.Errorf("error when dialing discord microservice: %s", err)
		}
	}

	return discord_provider.NewDiscordProviderClient(conn), nil
}

func initEmailProviderClient(ctx context.Context, v *viper.Viper, certFilePath string) (client email_provider.EmailProviderClient, err error) {
	hostname := v.GetString(config.FlagEmailProviderHostname)
	port := v.GetUint32(config.FlagEmailProviderPort)
	target := fmt.Sprintf("%s:%d", hostname, port)

	conn := &grpc.ClientConn{}
	if certFilePath != "" {
		creds, err := credentials.NewClientTLSFromFile(certFilePath, "")
		if err != nil {
			return nil, fmt.Errorf("error creating TLS credentials for the email microservice %v", err)
		}
		conn, err = grpc.Dial(target, grpc.WithTransportCredentials(creds))
		if err != nil {
			return nil, fmt.Errorf("error when dialing email microservice using TLS: %s", err)
		}
	} else {
		conn, err = grpc.Dial(target, grpc.WithInsecure())
		if err != nil {
			return nil, fmt.Errorf("error when dialing email microservice: %s", err)
		}
	}

	return email_provider.NewEmailProviderClient(conn), nil
}
