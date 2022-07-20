package config

import "github.com/spf13/cobra"

const (
	// General config settings
	FlagPort     = "port"
	FlagLogLevel = "loglevel"

	// Discord provider settings
	FlagDiscordProviderHostname = "discord_provider_hostname"
	FlagDiscordProviderPort     = "discord_provider_port"

	// Email provider settings
	FlagEmailProviderHostname = "email_provider_hostname"
	FlagEmailProviderPort     = "email_provider_port"

	// Database settings
	defaultDBPort uint32 = 5432
	dbPrefix             = "db-"
	FlagDbHost           = dbPrefix + "host"
	FlagDbPort           = dbPrefix + "port"
	FlagDbName           = dbPrefix + "name"
	FlagDbUser           = dbPrefix + "user"
	FlagDbSchema         = dbPrefix + "schema"
	FlagDbSslMode        = dbPrefix + "ssl"
)

func SetFlags(cmd *cobra.Command) {
	cmd.Flags().Uint32(FlagPort, 0, "The gRPC port on which to listen to")
	cmd.Flags().String(FlagLogLevel, "info", "The level of logs to print to stdout")

	cmd.Flags().String(FlagDiscordProviderHostname, "discord_provider", "The discord provider microservice's hostname to connect to")
	cmd.Flags().Uint32(FlagDiscordProviderPort, 0, "The discord provider microservice's port to connect to")

	cmd.Flags().String(FlagEmailProviderHostname, "email_provider", "The email provider microservice's hostname to connect to")
	cmd.Flags().Uint32(FlagEmailProviderPort, 0, "The email provider microservice's port to connect to")

	cmd.Flags().String(FlagDbHost, "postgres", "the database host to connect to")
	cmd.Flags().Uint32(FlagDbPort, defaultDBPort, "the database host to connect to")
	cmd.Flags().String(FlagDbName, "calidum_rotae", "the database name to connect to")
	cmd.Flags().String(FlagDbUser, "postgres", "the database user to use to connect to the database")
	cmd.Flags().String(FlagDbSslMode, "require", "[verify-ca, verify-full, require, disable]")
	cmd.Flags().String(FlagDbSchema, "calidum_rotae", "the database schema to use")

	cmd.MarkFlagRequired(FlagPort)
}
