package main

import (
	"context"

	"smart-choice/config"
	"smart-choice/database"
	"smart-choice/logger"
	"smart-choice/metrics"
	"smart-choice/middlewares"
	"smart-choice/routes"
	"smart-choice/tracing"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func main() {
	logger.InitLogger()

	tp, err := tracing.InitTracer()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize tracer")
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Error().Err(err).Msg("Error shutting down tracer provider")
		}
	}()

	config.LoadEnv()
	database.ConnectDB()

	r := gin.Default()

	r.Use(middlewares.SecurityHeadersMiddleware())
	r.Use(middlewares.CORSMiddleware())
	r.Use(middlewares.RateLimiter())
	r.Use(middlewares.LoggingMiddleware())
	r.Use(metrics.PrometheusMiddleware())
	r.Use(otelgin.Middleware("smart-choice"))

	routes.SetupRoutes(r)

	log.Info().Msg("Server starting on :8080")
	r.Run(":8080")
}
