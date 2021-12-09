package service

import (
	"fmt"

	"github.com/Weeping-Willow/api-example/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ScoreService interface {
	PostScore(score *models.RequestPostScore) (*models.DocumentScores, error)
}

type scoreService struct {
	commonServices
}

func newScoreService(opts *Options) *scoreService {
	return &scoreService{
		commonServices: commonServices{Repo: opts.Repo, Config: opts.Config},
	}
}

func (s *service) ScoreService() ScoreService {
	return s.scoreService
}

func (s *scoreService) PostScore(score *models.RequestPostScore) (*models.DocumentScores, error) {
	scoresFromDB, err := s.commonServices.Repo.GetScores(bson.M{"name": score.Name}, s.getDefaultOptions(1))
	if err != nil {
		return nil, err
	}

	finalScore, err := s.handlePostedScore(score, scoresFromDB)
	if err != nil {
		return nil, err
	}

	scoresWithRanks, err := s.getRanksForEachScore(finalScore)

	return scoresWithRanks[0], err
}

func (s *scoreService) getDefaultOptions(limit int64) *options.FindOptions {
	opts := options.Find()
	opts.SetLimit(limit)
	opts.SetSort(bson.M{"score": -1})
	return opts
}

func (s *scoreService) getEmptyRankMap(scores ...*models.DocumentScores) map[string]int {
	if len(scores) == 0 {
		return make(map[string]int)
	}
	rankings := make(map[string]int)
	for _, score := range scores {
		rankings[score.Name] = 0
	}

	return rankings
}

func (s *scoreService) handlePostedScore(postedScore *models.RequestPostScore, scoresFromDB []*models.DocumentScores) (*models.DocumentScores, error) {
	if len(scoresFromDB) > 0 {
		return s.handleExistingScoreUpdate(scoresFromDB[0], postedScore.Score)
	}

	endScore := &models.DocumentScores{
		Id:    primitive.NewObjectID(),
		Score: postedScore.Score,
		Name:  postedScore.Name,
	}
	_, err := s.Repo.InsertScore(endScore)
	if err != nil {
		return nil, err
	}

	return endScore, nil
}

func (s *scoreService) handleExistingScoreUpdate(score *models.DocumentScores, newScore int) (*models.DocumentScores, error) {
	if score.Score > newScore {
		return nil, fmt.Errorf(ErrScoreIsSmaller, newScore, score.Score)
	}
	score.Score = newScore

	updateResult, err := s.commonServices.Repo.UpdateScore(
		bson.M{"name": score.Name},
		score,
	)
	if err != nil {
		return nil, err
	}

	if updateResult.ModifiedCount == 0 {
		return nil, ErrUpdateFailed
	}

	return score, nil
}

func (s *scoreService) getRanksForEachScore(scores ...*models.DocumentScores) ([]*models.DocumentScores, error) {
	if len(scores) == 0 {
		return nil, nil
	}

	endRanks, err := s.Repo.GetScoreRanks(s.getEmptyRankMap(scores...), s.getDefaultOptions(0))
	if err != nil {
		return nil, err
	}

	for i, score := range scores {
		val, ok := endRanks[score.Name]
		if !ok {
			return nil, ErrRankingNotFound
		}
		scores[i].Rank = val
	}

	return scores, nil
}
