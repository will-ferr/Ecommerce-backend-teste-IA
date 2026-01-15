package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"smart-choice/config"
	"smart-choice/database"
	_ "smart-choice/docs"
	"smart-choice/logger"
	"smart-choice/metrics"
	"smart-choice/middlewares"
	"smart-choice/routes"
	"smart-choice/services"
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

	// Initialize Redis services
	serviceManager := services.GetServiceManager()
	if err := serviceManager.InitializeServices(); err != nil {
		log.Error().Err(err).Msg("Failed to initialize services")
	}

	r := gin.Default()

	r.Use(middlewares.SecurityHeadersMiddleware())
	r.Use(middlewares.CORSMiddleware())
	r.Use(middlewares.RateLimiter())
	r.Use(middlewares.LoggingMiddleware())
	r.Use(metrics.PrometheusMiddleware())
	r.Use(otelgin.Middleware("smart-choice"))

	routes.SetupRoutes(r)

	// Configure HTTP server with graceful shutdown
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Info().Msg("Server starting on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server failed to start")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")

	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
	}

	// Close database connection
	if sqlDB, err := database.DB.DB(); err == nil {
		sqlDB.Close()
		log.Info().Msg("Database connection closed")
	}

	// Shutdown services
	if err := serviceManager.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Error shutting down services")
	}

	log.Info().Msg("Server exited")
}
