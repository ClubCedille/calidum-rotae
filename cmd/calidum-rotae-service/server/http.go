package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/clubcedille/calidum-rotae-backend/cmd/calidum-rotae-service/config"
	"github.com/clubcedille/calidum-rotae-backend/cmd/calidum-rotae-service/instrumentation"
	"github.com/clubcedille/calidum-rotae-backend/pkg/calidum"
	"github.com/clubcedille/logger"
	serverutils "github.com/clubcedille/server-utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	ENV_CALIDUM_ROTAE_SERVICE_API_KEY = "CALIDUM_ROTAE_SERVICE_API_KEY"

	CALIDUM_ROTAE_TRACER_NAME = "calidum-rotae-tracer"

	EMAIL_POST_REQUEST   = "/email"
	DISCORD_POST_REQUEST = "/discord"
	DEFAULT_POST_REQUEST = "/"

	DISCORD_RPC_FUNC = "SendDiscordRpcRequest"
	EMAIL_RPC_FUNC   = "SendEmailRpcRequest"

	DISCORD_END_OF_SPAN = "Discord message sent!"
	EMAIL_END_OF_SPAN   = "Email sent!"
	OK_SPAN             = "HTTP request sent!"
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

	// Create the calidum rotae tracer
	calidumRotaeTracer := instrumentation.Traces{}
	calidumRotaeTracer.CalidumRotaeTracer = otel.Tracer(CALIDUM_ROTAE_TRACER_NAME)

	g.POST("/", func(g *gin.Context) { defaultPostRequest(g, services, calidumRotaeTracer) })
	g.POST("/discord", func(g *gin.Context) { discordPostRequest(g, services, calidumRotaeTracer) })
	g.POST("/email", func(g *gin.Context) { emailPostRequest(g, services, calidumRotaeTracer) })

	return g
}

func getRequestBody(g *gin.Context) ([]byte, error) {
	body, err := io.ReadAll(g.Request.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading http request")
	}
	return body, nil
}

func authenticationIsValid(g *gin.Context) bool {
	providedAPIKey := g.GetHeader("X-API-KEY")
	apiKey, found := os.LookupEnv(ENV_CALIDUM_ROTAE_SERVICE_API_KEY)
	return found && providedAPIKey == apiKey
}

func sendEmailRpcRequestWithSpan(ctx context.Context, g *gin.Context, body []byte, tracer instrumentation.Traces, services calidum.CalidumClient) (context.Context, trace.Span, error) {
	ctx, emailProviderGrpcSpan := tracer.GrpcSpan(ctx, EMAIL_RPC_FUNC, EMAIL_RPC_FUNC, instrumentation.EMAIL_PROVIDER_SERVICE)
	err := services.SendEmailRpcRequest(ctx, body)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		emailProviderGrpcSpan.RecordError(err)
		emailProviderGrpcSpan.SetStatus(codes.Error, err.Error())
	}

	return ctx, emailProviderGrpcSpan, err
}

func sendDiscordRpcRequestWithSpan(ctx context.Context, g *gin.Context, body []byte, tracer instrumentation.Traces, services calidum.CalidumClient) (context.Context, trace.Span, error) {
	ctx, discordGrpcSpan := tracer.GrpcSpan(ctx, DISCORD_RPC_FUNC, DISCORD_RPC_FUNC, instrumentation.DISCORD_PROVIDER_SERVICE)
	err := services.SendDiscordRpcRequest(ctx, body)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		discordGrpcSpan.RecordError(err)
		discordGrpcSpan.SetStatus(codes.Error, err.Error())
	}

	return ctx, discordGrpcSpan, err
}

// Send only an email
func emailPostRequest(g *gin.Context, services calidum.CalidumClient, tracer instrumentation.Traces) {
	ctx := g.Request.Context()

	ctx, httpSpan := tracer.HttpPostSpan(ctx, g, EMAIL_POST_REQUEST)
	defer httpSpan.End()

	if !authenticationIsValid(g) {
		g.JSON(http.StatusUnauthorized, gin.H{"error": "401 Unauthorized"})
		unauthorizedError := fmt.Errorf("error: 401 Unauthorized")
		httpSpan.RecordError(unauthorizedError)
		httpSpan.SetStatus(codes.Error, unauthorizedError.Error())
		return
	}

	body, err := getRequestBody(g)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		httpSpan.RecordError(err)
		httpSpan.SetStatus(codes.Error, err.Error())
		return
	}

	ctx, emailProviderGrpcSpan, err := sendEmailRpcRequestWithSpan(ctx, g, body, tracer, services)
	if err == nil {
		emailProviderGrpcSpan.SetStatus(codes.Ok, EMAIL_END_OF_SPAN)
		emailProviderGrpcSpan.End()
	}

	httpSpan.SetStatus(codes.Ok, OK_SPAN)
}

// Send only a discord message
func discordPostRequest(g *gin.Context, services calidum.CalidumClient, tracer instrumentation.Traces) {
	ctx := g.Request.Context()

	ctx, httpSpan := tracer.HttpPostSpan(ctx, g, DISCORD_POST_REQUEST)
	defer httpSpan.End()

	if !authenticationIsValid(g) {
		g.JSON(http.StatusUnauthorized, gin.H{"error": "401 Unauthorized"})
		unauthorizedError := fmt.Errorf("error: 401 Unauthorized")
		httpSpan.RecordError(unauthorizedError)
		httpSpan.SetStatus(codes.Error, unauthorizedError.Error())
		return
	}

	body, err := getRequestBody(g)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		httpSpan.RecordError(err)
		httpSpan.SetStatus(codes.Error, err.Error())
		return
	}

	ctx, discordGrpcSpan, err := sendDiscordRpcRequestWithSpan(ctx, g, body, tracer, services)
	if err == nil {
		discordGrpcSpan.SetStatus(codes.Ok, DISCORD_END_OF_SPAN)
		discordGrpcSpan.End()
	}

	httpSpan.SetStatus(codes.Ok, OK_SPAN)
}

// Send an email and a discord message
func defaultPostRequest(g *gin.Context, services calidum.CalidumClient, tracer instrumentation.Traces) {
	ctx := g.Request.Context()

	ctx, httpSpan := tracer.HttpPostSpan(ctx, g, DEFAULT_POST_REQUEST)
	defer httpSpan.End()

	if !authenticationIsValid(g) {
		g.JSON(http.StatusUnauthorized, gin.H{"error": "401 Unauthorized"})
		unauthorizedError := fmt.Errorf("error: 401 Unauthorized")
		httpSpan.RecordError(unauthorizedError)
		httpSpan.SetStatus(codes.Error, unauthorizedError.Error())
		return
	}

	body, err := getRequestBody(g)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		httpSpan.RecordError(err)
		httpSpan.SetStatus(codes.Error, err.Error())
		return
	}

	_, discordGrpcSpan, err := sendDiscordRpcRequestWithSpan(ctx, g, body, tracer, services)
	if err == nil {
		discordGrpcSpan.SetStatus(codes.Ok, DISCORD_END_OF_SPAN)
		discordGrpcSpan.End()
	}

	_, emailProviderGrpcSpan, err := sendEmailRpcRequestWithSpan(ctx, g, body, tracer, services)
	if err == nil {
		emailProviderGrpcSpan.SetStatus(codes.Ok, EMAIL_END_OF_SPAN)
		emailProviderGrpcSpan.End()
	}

	httpSpan.SetStatus(codes.Ok, OK_SPAN)
}
