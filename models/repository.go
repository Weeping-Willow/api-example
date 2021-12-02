package models

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
}

type repo struct {
	db *mongo.Database
}

const CollectionNameScores = "scores"

func GetConnection(dbHost string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(dbHost)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to ping database at %s with error: %w", dbHost, err)
	}
	return client, nil
}

func NewRepository(database *mongo.Database) Repository {
	return &repo{
		db: database,
	}
}

func (r *repo) collection(name string) *mongo.Collection {
	return r.db.Collection(name)
}
