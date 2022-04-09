package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/clubcedille/calidum-rotae-backend/cmd/calidum-rotae-service/app"
	"github.com/clubcedille/calidum-rotae-backend/cmd/calidum-rotae-service/config"
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

	// Start the microservice service and its dependencies.
	if _, err := app.Run(ctx, v); err != nil {
		return fmt.Errorf("error when initializing calidum rotae service: %s", err)
	}

	return nil
}

func main() {
	if err := serviceCmd.Execute(); err != nil {
		log.Fatalf("error when running the calidum rotae service: %s\n", err)
		os.Exit(1)
	}
}
