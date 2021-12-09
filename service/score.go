package service

import (
	"context"
	"fmt"

	"github.com/Weeping-Willow/api-example/models"
	"github.com/Weeping-Willow/api-example/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ScoreService interface {
	PostScore(score *models.RequestPostScore) (*models.DocumentScores, error)
	GetScores(req *models.RequestGetScores) (*models.ScoreResponse, error)
}

type scoreService struct {
	commonServices
	ScorePagination *scorePagination
}

type scorePagination struct {
	PerPage     int
	CurrentPage int
	MaxPage     int
	NextPage    int
	Total       int
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

func (s *scoreService) GetScores(req *models.RequestGetScores) (*models.ScoreResponse, error) {
	total, err := s.Repo.Collection(repositories.CollectionNameScores).CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}

	s.ScorePagination = &scorePagination{}
	if err := s.ScorePagination.setUpPagination(total, req.PageSize, req.PageNumber); err != nil {
		return nil, err
	}

	opts := s.getDefaultOptions(int64(s.ScorePagination.PerPage))
	if s.ScorePagination.CurrentPage > 1 {
		opts.SetSkip(int64((s.ScorePagination.CurrentPage - 1) * s.ScorePagination.PerPage))
	}

	scores, err := s.Repo.GetScores(bson.M{}, opts)
	if err != nil {
		return nil, err
	}

	scores, err = s.getRanksForEachScore(scores...)
	if err != nil {
		return nil, err
	}

	if s.ScorePagination.CurrentPage < s.ScorePagination.MaxPage {
		s.ScorePagination.NextPage = s.ScorePagination.CurrentPage + 1
	}

	var aroundMeScores []*models.DocumentScores
	aroundMeScores, err = s.handleNameFilter(scores, req.Name)
	if err != nil {
		return nil, err
	}

	return &models.ScoreResponse{
		Results: scores,
		Meta: &models.Meta{
			CurrentPage: s.ScorePagination.CurrentPage,
			TotalPages:  s.ScorePagination.MaxPage,
			TotalCount:  s.ScorePagination.Total,
			NextPage:    s.ScorePagination.NextPage,
			PerPage:     s.ScorePagination.PerPage,
		},
		AroundMe: aroundMeScores,
	}, err
}

func (s *scorePagination) setUpPagination(total int64, perPage, currentPage int) error {
	if perPage < 0 || currentPage < 0 {
		return ErrPaginationInvalid
	}
	if total == 0 {
		return mongo.ErrNoDocuments
	}

	s.PerPage = perPage
	s.CurrentPage = currentPage
	s.Total = int(total)

	if s.PerPage == 0 {
		s.PerPage = 10
	}
	if currentPage == 0 {
		s.CurrentPage = 1
	}
	s.MaxPage = s.Total / s.PerPage
	if s.Total%s.PerPage > 0 {
		s.MaxPage++
	}

	if s.CurrentPage > s.MaxPage {
		return ErrPaginationPageMoreThanMax
	}

	return nil
}

func (s *scoreService) handleNameFilter(scores []*models.DocumentScores, name string) ([]*models.DocumentScores, error) {
	if name == "" {
		return nil, nil
	}

	for _, score := range scores {
		if score.Name == name {
			return nil, nil
		}
	}

	s.commonServices.Repo.GetScoreAndNeighbors(name, s.getDefaultOptions(0))
	return s.commonServices.Repo.GetScoreAndNeighbors(name, s.getDefaultOptions(0))
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
	if score.Score >= newScore {
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

func (s *scoreService) getDefaultOptions(limit int64) *options.FindOptions {
	opts := options.Find()
	opts.SetLimit(limit)
	opts.SetSort(bson.M{"score": -1})
	opts.SetSkip(0)
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
