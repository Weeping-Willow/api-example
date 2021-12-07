package main

import (
	"context"
	"os"

	"github.com/Weeping-Willow/api-example/config"
	"github.com/Weeping-Willow/api-example/repositories"
	"github.com/Weeping-Willow/api-example/router"
	"github.com/Weeping-Willow/api-example/service"
	"github.com/rs/zerolog/log"
)

func main() {
	if err := run(); err != nil {
		log.Fatal().Err(err).Msg("failed to run")
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.New()
	if err != nil {
		return err
	}

	client, err := repositories.GetConnection(cfg.DatabaseUrl)
	if err != nil {
		return err
	}
	defer client.Disconnect(context.TODO())
	repo := repositories.NewRepository(client.Database(cfg.DatabaseName))

	s := service.NewService(&service.Options{
		Repo:   repo,
		Config: cfg,
	})

	log.Info().Msg("starting server")
	return router.StartServer(s)
}
