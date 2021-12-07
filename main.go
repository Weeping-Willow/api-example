package main

import (
	"context"
	"math/rand"
	"os"

	"github.com/Weeping-Willow/api-example/config"
	"github.com/Weeping-Willow/api-example/models"
	"github.com/Weeping-Willow/api-example/repositories"
	"github.com/Weeping-Willow/api-example/router"
	"github.com/Weeping-Willow/api-example/service"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	seedDb(repo)

	s := service.NewService(&service.Options{
		Repo:   repo,
		Config: cfg,
	})

	log.Info().Msg("starting server")
	return router.StartServer(s)
}

func seedDb(repo repositories.MongoRepository) {
	repo.DeleteAll(repositories.CollectionNameScores)
	repo.Collection(repositories.CollectionNameScores).Indexes().CreateOne(
		context.TODO(),
		mongo.IndexModel{
			Keys: bson.D{{Key: "name", Value: 1}},
		},
	)
	repo.Collection(repositories.CollectionNameScores).Indexes().CreateOne(
		context.TODO(),
		mongo.IndexModel{
			Keys: bson.D{{Key: "score", Value: 1}},
		},
	)

	for i := 0; i < 50; i++ {
		repo.InsertScore(&models.DocumentScores{
			Id:    primitive.NewObjectID(),
			Name:  randSeq(4),
			Score: rand.Intn(500),
		})
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
