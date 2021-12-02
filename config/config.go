package config

import (
	"net/url"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Port         string
	DatabaseUrl  string
	DatabaseName string
}

func New(filenames ...string) (*Config, error) {
	if err := godotenv.Load(filenames...); err != nil {
		return nil, err
	}
	return &Config{
		Port:         getEnvWithDefault("PORT", "8080"),
		DatabaseUrl:  getEnvWithDefault("DATABASE_URL", ""),
		DatabaseName: getDatabaseNameFromMongoURL(getEnvWithDefault("DATABASE_URL", "")),
	}, nil
}

func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getDatabaseNameFromMongoURL(dbURL string) string {
	u, err := url.Parse(dbURL)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to parse database url")
		return ""
	}
	database, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to parse database url query")
		return ""
	}

	if database["authSource"] == nil {
		log.Warn().Msg("Failed to find authSource in database url query")
		return ""
	}

	return database["authSource"][0]
}
