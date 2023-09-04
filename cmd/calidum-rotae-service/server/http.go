package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/clubcedille/calidum-rotae-backend/cmd/calidum-rotae-service/config"
	"github.com/clubcedille/calidum-rotae-backend/pkg/calidum"
	"github.com/clubcedille/logger"
	serverutils "github.com/clubcedille/server-utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

const (
	ENV_CALIDUM_ROTAE_SERVICE_API_KEY = "CALIDUM_ROTAE_SERVICE_API_KEY"
)

func InitHTTPServerFromViper(ctx context.Context, v *viper.Viper, services calidum.CalidumClient) (*serverutils.HttpServer, error) {
	addr := fmt.Sprintf(":%d", v.GetUint32(config.FlagPort))

	// HTTP Server configuration
	httpServer := &http.Server{}
	httpServer.Addr = addr
	httpServer.Handler = initHTTPServerHandler(ctx, services)

	return serverutils.NewHttpServer(httpServer), nil
}

func initHTTPServerHandler(ctx context.Context, services calidum.CalidumClient) *gin.Engine {
	g := gin.New()
	ctxLogger := logger.NewFromContextOrDefault(ctx)

	// Middlewares configuration
	g.Use(logger.HTTPLoggerMiddleware(ctxLogger))

	// TODO: Database http request
	g.POST("/", func(g *gin.Context) { defaultPostRequest(g, services) })
	g.POST("/discord", func(g *gin.Context) { discordPostRequest(g, services) })
	g.POST("/email", func(g *gin.Context) { emailPostRequest(g, services) })

	return g
}

func authenticationIsValid(g *gin.Context) bool {
	providedAPIKey := g.GetHeader("X-API-KEY")
	apiKey, found := os.LookupEnv(ENV_CALIDUM_ROTAE_SERVICE_API_KEY)
	return found && providedAPIKey == apiKey
}

func getRequestBody(g *gin.Context) ([]byte, error) {
	body, err := io.ReadAll(g.Request.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading http request")
	}
	return body, nil
}

func emailPostRequest(g *gin.Context, services calidum.CalidumClient) {
	if !authenticationIsValid(g) {
		g.JSON(http.StatusUnauthorized, gin.H{"error": " 401 Unauthorized"})
		return
	}

	body, err := getRequestBody(g)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = services.SendEmailRpcRequest(g.Request.Context(), body)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func discordPostRequest(g *gin.Context, services calidum.CalidumClient) {
	if !authenticationIsValid(g) {
		g.JSON(http.StatusUnauthorized, gin.H{"error": " 401 Unauthorized"})
		return
	}

	body, err := getRequestBody(g)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = services.SendDiscordRpcRequest(g.Request.Context(), body)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

// Call all microservices.
func defaultPostRequest(g *gin.Context, services calidum.CalidumClient) {
	if !authenticationIsValid(g) {
		g.JSON(http.StatusUnauthorized, gin.H{"error": " 401 Unauthorized"})
		return
	}

	body, err := getRequestBody(g)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = services.SendDiscordRpcRequest(g.Request.Context(), body)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	err = services.SendEmailRpcRequest(g.Request.Context(), body)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
