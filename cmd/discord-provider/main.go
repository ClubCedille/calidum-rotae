package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/clubcedille/calidum-rotae-backend/cmd/discord-provider/config"
	"github.com/clubcedille/calidum-rotae-backend/cmd/discord-provider/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	discordCmd = &cobra.Command{
		Use:          "discord_provider",
		Short:        "Discord provider microservice for the calidum rotae app",
		SilenceUsage: true,
		RunE:         runDiscordProvider,
	}
)

func init() {
	discordCmd.Flags().Uint32(config.FlagPort, 0, "The gRPC port on which to listen to")
	discordCmd.Flags().String(config.FlagLogLevel, "info", "The level of logs to print to stdout")

	discordCmd.MarkFlagRequired(config.FlagPort)
}

func runDiscordProvider(cmd *cobra.Command, args []string) (err error) {
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
	if err := discordCmd.Execute(); err != nil {
		log.Fatalf("error when running the discord provider: %s\n", err)
		os.Exit(1)
	}
}
