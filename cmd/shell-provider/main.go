package main

import (
	"context"
    "fmt"
    "log"
    "os"

	"github.com/clubcedille/calidum-rotae-backend/cmd/shell-provider/config"
	"github.com/clubcedille/calidum-rotae-backend/cmd/shell-provider/server"
	"github.com/clubcedille/logger"
	serverutils "github.com/clubcedille/server-utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	shellCmd = &cobra.Command{
		Use:          "shell_provider",
		Short:        "Shell provider microservice for the calidum rotae app",
		SilenceUsage: true,
		RunE:         runShellProvider,
	}
)

func init() {
	shellCmd.Flags().Uint32(config.FlagPort, 0, "The gRPC port on which to listen to")
	shellCmd.Flags().String(config.FlagLogLevel, "info", "The level of logs to print to stdout")

	shellCmd.MarkFlagRequired(config.FlagPort)
}

func runShellProvider(cmd *cobra.Command, args []string) (err error) {
	v := viper.New()
	if err := v.BindPFlags(cmd.Flags()); err != nil {
		return fmt.Errorf("error when binding flags: %s", err)
	}

	// Initialize context instance with logger
	ctxLogger := logger.Initialize(logger.Config{
		Level: v.GetString(config.FlagLogLevel),
	})
	ctx := context.WithValue(cmd.Context(), logger.CtxKey, ctxLogger)

	// Initialize server
	v.AutomaticEnv()
	port := v.GetUint32(config.FlagPort)
	server := server.NewServer()
	grpcServer := server.ConfigureGrpc()
	grpcServerInstance := serverutils.NewGrpcServer(grpcServer)

	if err := grpcServerInstance.Run(ctx, serverutils.RunRequest{
		Port:              int32(port),
		ShutdownTimeoutMs: 10000, // 10 seconds
	}); err != nil {
		return fmt.Errorf("error serving gRPC server over port %d: %s", port, err)
	}

	return
}

func main() {
	if err := shellCmd.Execute(); err != nil {
		log.Fatalf("error when running the shell provider: %s\n", err)
		os.Exit(1)
	}
}
