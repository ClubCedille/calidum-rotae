package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/clubcedille/calidum-rotae-backend/cmd/email-provider/config"
	"github.com/clubcedille/calidum-rotae-backend/cmd/email-provider/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	emailCmd = &cobra.Command{
		Use:          "email_provider",
		Short:        "Email provider microservice for the calidum rotae app",
		SilenceUsage: true,
		RunE:         runEmailProvider,
	}
)

func init() {
	emailCmd.Flags().Uint32(config.FlagPort, 0, "The gRPC port on which to listen to")
	emailCmd.Flags().String(config.FlagLogLevel, "info", "The level of logs to print to stdout")

	emailCmd.MarkFlagRequired(config.FlagPort)
}

func runEmailProvider(cmd *cobra.Command, args []string) (err error) {
	v := viper.New()
	if err := v.BindPFlags(cmd.Flags()); err != nil {
		return fmt.Errorf("error when binding flags: %s", err)
	}

	// TODO: Maybe use this ctx in the future.
	// ctxLogger := logger.Initialize(logger.Config{
	// 	Level: v.GetString(config.FlagLogLevel),
	// })
	// ctx := context.WithValue(cmd.Context(), logger.CtxKey, ctxLogger)

	v.AutomaticEnv()
	port := v.GetUint32(config.FlagPort)
	lt, err := net.Listen("tcp", fmt.Sprintf(":%d", v.GetUint32(config.FlagPort)))
	if err != nil {
		return fmt.Errorf("error starting tcp listener on port %d: %s", port, err)
	}

	server := server.NewServer()
	grpcServer := server.ConfigureGrpc()
	if err := grpcServer.Serve(lt); err != nil {
		return fmt.Errorf("failed to serve gRPC server over port %d: %s", port, err)
	}

	return
}

func main() {
	if err := emailCmd.Execute(); err != nil {
		log.Fatalf("error when running the email provider: %s\n", err)
		os.Exit(1)
	}
}
