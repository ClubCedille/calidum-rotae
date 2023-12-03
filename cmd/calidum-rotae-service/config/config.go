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

	// OTEL Tracer settings
	FlagOTELOtlpExporterHost = "otel_otlp_exporter_host"
	FlagOTELOtlpExporterPort = "otel_otlp_exporter_port"
)

func SetFlags(cmd *cobra.Command) {
	cmd.Flags().Uint32(FlagPort, 0, "The gRPC port on which to listen to")
	cmd.Flags().String(FlagLogLevel, "info", "The level of logs to print to stdout")

	cmd.Flags().String(FlagDiscordProviderHostname, "discord_provider", "The discord provider microservice's hostname to connect to")
	cmd.Flags().Uint32(FlagDiscordProviderPort, 0, "The discord provider microservice's port to connect to")

	cmd.Flags().String(FlagEmailProviderHostname, "email_provider", "The email provider microservice's hostname to connect to")
	cmd.Flags().Uint32(FlagEmailProviderPort, 0, "The email provider microservice's port to connect to")

	cmd.Flags().String(FlagOTELOtlpExporterHost, "localhost", "The OTEL OTLP exporter host to connect to")
	cmd.Flags().String(FlagOTELOtlpExporterPort, "4318", "The OTEL OTLP exporter port to connect to")

	cmd.MarkFlagRequired(FlagPort)
}
