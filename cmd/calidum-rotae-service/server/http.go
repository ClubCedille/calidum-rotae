package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/clubcedille/calidum-rotae-backend/cmd/calidum-rotae-service/config"
	"github.com/clubcedille/logger"
	serverutils "github.com/clubcedille/server-utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func InitHTTPServerFromViper(ctx context.Context, v *viper.Viper) (*serverutils.HttpServer, error) {
	addr := fmt.Sprintf(":%d", v.GetUint32(config.FlagPort))

	// HTTP Server configuration
	httpServer := &http.Server{}
	httpServer.Addr = addr
	httpServer.Handler = initHTTPServerHandler(ctx)

	return serverutils.NewHttpServer(httpServer), nil
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
