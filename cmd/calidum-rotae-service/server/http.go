package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/clubcedille/calidum-rotae-backend/cmd/calidum-rotae-service/config"
	"github.com/clubcedille/logger"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type HTTPServer struct {
	server *http.Server
}

func InitHTTPServerFromViper(ctx context.Context, v *viper.Viper) (*HTTPServer, error) {
	addr := fmt.Sprintf(":%d", v.GetUint32(config.FlagPort))

	// HTTP Server configuration
	httpServer := &http.Server{}
	httpServer.Addr = addr
	httpServer.Handler = initHTTPServerHandler(ctx)

	return &HTTPServer{httpServer}, nil
}

func (s HTTPServer) Serve() error {
	return s.server.ListenAndServe()
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
