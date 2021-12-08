package repositories

import (
	"context"
	"fmt"

	"github.com/Weeping-Willow/api-example/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const CollectionNameScores = "scores"

type MongoRepository interface {
	InsertScore(score *models.DocumentScores) (*mongo.InsertOneResult, error)
	GetScores(filters primitive.M, opts ...*options.FindOptions) ([]*models.DocumentScores, error)
	UpdateScore(filter primitive.M, score *models.DocumentScores) (*mongo.UpdateResult, error)
	DeleteAll(collectionName string) (*mongo.DeleteResult, error)
	Collection(name string) *mongo.Collection
	GetScoreRanks(scores map[string]int, opts ...*options.FindOptions) (map[string]int, error)
}

type repo struct {
	db *mongo.Database
}

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

func NewRepository(database *mongo.Database) MongoRepository {
	return &repo{
		db: database,
	}
}

func (r *repo) Collection(name string) *mongo.Collection {
	return r.db.Collection(name)
}

func (r *repo) DeleteAll(collectionName string) (*mongo.DeleteResult, error) {
	collection := r.Collection(collectionName)
	res, err := collection.DeleteMany(context.TODO(), bson.M{})
	return res, err
}
