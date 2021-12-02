package service

import (
	"github.com/Weeping-Willow/api-example/config"
	"github.com/Weeping-Willow/api-example/models"
)

type Service interface {
	GetConfig() *config.Config
}

type Options struct {
	Repo   models.Repository
	Config *config.Config
}

type service struct {
	Options *Options
}

func NewService(o *Options) Service {
	return &service{
		Options: o,
	}
}

func (s *service) GetConfig() *config.Config {
	return s.Options.Config
}
