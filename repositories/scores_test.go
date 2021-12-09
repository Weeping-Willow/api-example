package repositories

import (
	"fmt"
	"testing"

	"github.com/Weeping-Willow/api-example/models"
	"github.com/Weeping-Willow/api-example/testt"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Test_repo_GetScores(t *testing.T) {
	db := testt.DBWithModelMocks(t,
		func() interface{} {
			tests := make([]models.DocumentScores, 10)
			for i := 0; i < 10; i++ {
				tests[i] = models.DocumentScores{Id: primitive.NewObjectID(), Score: i, Name: fmt.Sprintf("name%d", i)}
			}
			return tests
		},
		CollectionNameScores,
	)

	tests := []struct {
		name      string
		giveCount int
		wantCount int
	}{
		{
			name:      "Give zero and get all",
			giveCount: 0,
			wantCount: 10,
		},
		{
			name:      "Give 5 and get 5",
			giveCount: 5,
			wantCount: 5,
		},
		{
			name:      "Give 10 and get 10",
			giveCount: 10,
			wantCount: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRepository(db)
			gotResults, gotErr := r.GetScores(primitive.M{}, options.Find().SetLimit(int64(tt.giveCount)))
			assert.Nil(t, gotErr)
			assert.Len(t, gotResults, tt.wantCount)
		})
	}
}

func Test_repo_GetScoreRanks(t *testing.T) {
	db := testt.DBWithModelMocks(t,
		func() interface{} {
			tests := make([]models.DocumentScores, 4)
			for i := 0; i < 4; i++ {
				tests[i] = models.DocumentScores{Id: primitive.NewObjectID(), Score: i, Name: fmt.Sprintf("name%d", i)}
			}
			return tests
		},
		CollectionNameScores,
	)

	tests := []struct {
		name      string
		giveRanks map[string]int
		wantRanks map[string]int
	}{
		{
			name: "Give all ranks and all get ranks",
			giveRanks: map[string]int{
				"name0": 0,
				"name1": 0,
				"name2": 0,
				"name3": 0,
			},
			wantRanks: map[string]int{
				"name0": 4,
				"name1": 3,
				"name2": 2,
				"name3": 1,
			},
		},
		{
			name: "Give one score and get rank for it",
			giveRanks: map[string]int{
				"name2": 0,
			},
			wantRanks: map[string]int{
				"name2": 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRepository(db)
			gotResults, gotErr := r.GetScoreRanks(tt.giveRanks, options.Find().SetSort(primitive.M{"score": -1}))
			assert.Nil(t, gotErr)
			assert.Equal(t, gotResults, tt.wantRanks)
		})
	}
}
