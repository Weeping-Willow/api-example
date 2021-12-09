package service

import (
	"fmt"
	"testing"

	"github.com/Weeping-Willow/api-example/models"
	"github.com/Weeping-Willow/api-example/repositories"
	"github.com/Weeping-Willow/api-example/testt"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Test_scoreService_PostScore(t *testing.T) {
	db := testt.DBWithModelMocks(t,
		func() interface{} {
			tests := make([]models.DocumentScores, 4)
			for i := 0; i < 4; i++ {
				tests[i] = models.DocumentScores{Id: primitive.NewObjectID(), Score: (i + 1) * 2, Name: fmt.Sprintf("name%d", i)}
			}
			return tests
		},
		repositories.CollectionNameScores,
	)

	tests := []struct {
		name        string
		giveRequest *models.RequestPostScore
		want        *models.DocumentScores
		wantErr     error
	}{
		{
			name: "Success: add new score and return it and its rank",
			giveRequest: &models.RequestPostScore{
				Name:  "test",
				Score: 10,
			},
			want: &models.DocumentScores{
				Score: 10,
				Name:  "test",
				Rank:  1,
			},
		},
		{
			name: "Success: add existing score that is bigger, return it and its rank",
			giveRequest: &models.RequestPostScore{
				Name:  "name1",
				Score: 7,
			},
			want: &models.DocumentScores{
				Name:  "name1",
				Score: 7,
				Rank:  3,
			},
		},
		{
			name: "Failure: post requests with smaller than current score",
			giveRequest: &models.RequestPostScore{
				Name:  "name1",
				Score: 3,
			},
			want:    nil,
			wantErr: fmt.Errorf(ErrScoreIsSmaller, 3, 7),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &scoreService{
				commonServices: commonServices{
					Repo: repositories.NewRepository(db),
				},
			}
			got, err := s.PostScore(tt.giveRequest)
			if err == nil {
				tt.want.Id = got.Id
			}
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
