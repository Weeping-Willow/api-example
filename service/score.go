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
	scores, err := s.commonServices.Repo.GetScores(bson.M{"name": score.Name}, s.getDefaultOptions(1))
	if err != nil {
		return nil, err
	}

	//TODO: The if should be simplified in some way
	var endScore *models.DocumentScores
	if len(scores) > 0 {
		endScore = scores[0]
		if endScore.Score < score.Score {
			return nil, fmt.Errorf("given score %d is smaller than already existing score %d", score.Score, scores[0].Score)
		}
		updateResult, err := s.commonServices.Repo.UpdateScore(
			bson.M{"name": score.Name},
			&models.DocumentScores{
				Id:    endScore.Id,
				Score: score.Score,
				Name:  score.Name,
			},
		)
		if err != nil {
			return nil, err
		}

		if updateResult.ModifiedCount == 0 {
			return nil, ErrUpdateFailed
		}
	} else {
		endScore = &models.DocumentScores{
			Id:    primitive.NewObjectID(),
			Score: score.Score,
			Name:  score.Name,
		}
		_, err := s.Repo.InsertScore(endScore)
		if err != nil {
			return nil, err
		}
	}

	// TODO refactor
	endRanks, err := s.Repo.GetScoreRanks(s.getEmptyRankMap(endScore), s.getDefaultOptions(0))
	if err != nil {
		return nil, err
	}
	val, ok := endRanks[endScore.Name]
	if !ok {
		return nil, ErrRankingNotFound
	}
	endScore.Rank = val

	return endScore, err
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
