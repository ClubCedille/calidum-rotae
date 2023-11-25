package main

import (
	"context"
	"fmt"
	"log"

	"github.com/clubcedille/calidum-rotae-backend/cmd/calidum-rotae-service/app"
	"github.com/clubcedille/calidum-rotae-backend/cmd/calidum-rotae-service/config"
	instrumentation "github.com/clubcedille/calidum-rotae-backend/cmd/otel-instrumentation"
	"github.com/clubcedille/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	serviceCmd = &cobra.Command{
		Use:          "calidum_rotae",
		Short:        "Calidum Rotae microservice for the calidum rotae app",
		SilenceUsage: true,
		RunE:         runService,
	}
)

func init() {
	config.SetFlags(serviceCmd)
}

func runService(cmd *cobra.Command, args []string) error {
	v := viper.New()
	if err := v.BindPFlags(cmd.Flags()); err != nil {
		return fmt.Errorf("error when binding flags: %s", err)
	}

	// Setup context + logger
	ctxLogger := logger.Initialize(logger.Config{
		Level: v.GetString(config.FlagLogLevel),
	})
	ctx := context.WithValue(cmd.Context(), logger.CtxKey, ctxLogger)

	// Initialize the service and its dependencies
	service, err := app.InitFromViper(ctx, v)
	if err != nil {
		return fmt.Errorf("error when initializing calidum rotae service: %s", err)
	}

	// Setup OpenTelemetry
	otlpHost := v.GetString(config.FlagOTELOtlpExporterHost)
	otlpPort := v.GetString(config.FlagOTELOtlpExporterPort)
	tp, err := instrumentation.SetupOpenTelemetry(ctx, otlpHost, otlpPort)
	if err != nil {
		log.Fatalf("error when setting up OpenTelemetry: %s\n", err)
	}

	defer func() { _ = tp.Shutdown(ctx) }()

	// Start the microservice service and its dependencies.
	if err := service.Run(
		ctx,
		int32(v.GetUint32(config.FlagPort)),
	); err != nil {
		return fmt.Errorf("error when running calidum rotae: %s", err)
	}

	return nil
}

func main() {
	if err := serviceCmd.Execute(); err != nil {
		log.Fatalf("error when running the calidum rotae service: %s\n", err)
	}
}
