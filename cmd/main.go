package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/infrastructure/config"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/infrastructure/di"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/infrastructure/inbound/rest/routes"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/logger"
)

func main() {
	cfg := config.NewConfig()

	if cfg.GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	app, err := di.InitializeMediator()

	if err != nil {
		log.Fatal().Msgf("Error executing DI: %v", err)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(logger.GinZerologMiddleware())

	routes.RegisterRoutes(app, router)

	port := fmt.Sprintf(":%s", cfg.Port)

	log.Info().
		Str("port", port).
		Str("environment", cfg.AppEnv).
		Msg("Server running")

	if err := router.Run(port); err != nil {
		log.Fatal().
			Err(err).
			Msg("failed to start server")
	}
}
