package service

import (
	"github.com/Weeping-Willow/api-example/config"
	"github.com/Weeping-Willow/api-example/repositories"
)

type Service interface {
	GetConfig() *config.Config
	TokenService() TokenService
	ScoreService() ScoreService
}

type Options struct {
	Repo   repositories.MongoRepository
	Config *config.Config
}

type service struct {
	Options       *Options
	tokensService TokenService
	scoreService  ScoreService
}

type commonServices struct {
	Repo   repositories.MongoRepository
	Config *config.Config
}

func NewService(o *Options) Service {
	return &service{
		Options:       o,
		tokensService: newTokenService(o),
		scoreService:  newScoreService(o),
	}
}

func (s *service) GetConfig() *config.Config {
	return s.Options.Config
}
