package main

import (
	"os"

	"github.com/Weeping-Willow/api-example/config"
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

	log.Info().Msg(cfg.Port)

	return nil
}
