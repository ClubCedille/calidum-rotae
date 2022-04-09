package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/clubcedille/calidum-rotae-backend/cmd/calidum-rotae-service/config"
	"github.com/clubcedille/logger"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func InitHTTPServerFromViper(ctx context.Context, v *viper.Viper) (*http.Server, error) {
	addr := fmt.Sprintf(":%d", viper.GetUint32(config.FlagPort))

	// HTTP Server configuration
	httpServer := &http.Server{}
	httpServer.Addr = addr
	httpServer.Handler = initHTTPServerHandler(ctx)

	// Make sure socket can be opened
	_, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("error opening socket on address %s: %s", addr, err)
	}

	return httpServer, nil
}

func initHTTPServerHandler(ctx context.Context) *gin.Engine {
	g := gin.New()
	ctxLogger := logger.NewFromContextOrDefault(ctx)

	// Middlewares configuration
	g.Use(logger.HTTPLoggerMiddleware(ctxLogger))

	// TODO: Define all HTTP routes here
	// g.GET("/...")

	return g
}
