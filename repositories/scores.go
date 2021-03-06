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

func (r *repo) InsertScore(score *models.DocumentScores) (*mongo.InsertOneResult, error) {
	res, err := r.Collection(CollectionNameScores).InsertOne(context.TODO(), score)
	return res, err
}

func (r *repo) GetScores(filter primitive.M, opts ...*options.FindOptions) ([]*models.DocumentScores, error) {
	cursor, err := r.Collection(CollectionNameScores).Find(context.TODO(), filter, opts...)
	if err != nil {
		return nil, fmt.Errorf("error finding scores: %w", err)
	}
	defer cursor.Close(context.TODO())

	var scores []*models.DocumentScores
	for cursor.Next(context.TODO()) {
		var score *models.DocumentScores
		if err := cursor.Decode(&score); err != nil {
			return scores, fmt.Errorf("error decoding score: %w", err)
		}
		scores = append(scores, score)
	}
	return scores, nil
}

func (r *repo) GetScoreRanks(scores map[string]int, opts ...*options.FindOptions) (map[string]int, error) {
	maxScores := len(scores)
	scoresFound := 0
	i := 1

	cursor, err := r.Collection(CollectionNameScores).Find(context.TODO(), bson.M{}, opts...)
	if err != nil {
		return nil, fmt.Errorf("error finding scores: %w", err)
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		if maxScores == scoresFound {
			break
		}
		var score *models.DocumentScores
		if err := cursor.Decode(&score); err != nil {
			return scores, fmt.Errorf("error decoding score: %w", err)
		}
		if _, ok := scores[score.Name]; ok {
			maxScores++
			scores[score.Name] = i
		}
		i++
	}
	return scores, nil
}

func (r *repo) GetScoreAndNeighbors(name string, opts ...*options.FindOptions) ([]*models.DocumentScores, error) {
	var previousScore *models.DocumentScores
	scores := []*models.DocumentScores{}
	addNextScore := false
	i := 1

	cursor, err := r.Collection(CollectionNameScores).Find(context.TODO(), bson.M{}, opts...)
	if err != nil {
		return nil, fmt.Errorf("error finding scores: %w", err)
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var score *models.DocumentScores
		if err := cursor.Decode(&score); err != nil {
			return scores, fmt.Errorf("error decoding score: %w", err)
		}
		score.Rank = i
		if addNextScore {
			scores = append(scores, score)
			break
		}

		if score.Name == name {
			if previousScore != nil {
				scores = append(scores, previousScore)
			}
			scores = append(scores, score)
			addNextScore = true
		}
		previousScore = score
		i++
	}
	return scores, nil
}

func (r *repo) UpdateScore(filter primitive.M, score *models.DocumentScores) (*mongo.UpdateResult, error) {
	res, err := r.Collection(CollectionNameScores).ReplaceOne(context.TODO(), filter, score)
	return res, err
}
