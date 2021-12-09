package router

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/Weeping-Willow/api-example/models"
	"github.com/Weeping-Willow/api-example/repositories"
	"github.com/Weeping-Willow/api-example/service"
	"github.com/Weeping-Willow/api-example/testt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Test_controllerPostScore(t *testing.T) {
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
		name       string
		wantBody   string
		wantStatus int
		giveBody   string
	}{
		{
			name:       "Failure: Empty body",
			giveBody:   ``,
			wantBody:   `{"error":"EOF"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Failure: no name",
			giveBody:   `{"score": 3}`,
			wantBody:   `{"error":"Key: 'RequestPostScore.Name' Error:Field validation for 'Name' failed on the 'required' tag"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Failure: no score",
			giveBody:   `{"name": "test"}`,
			wantBody:   `{"error":"Key: 'RequestPostScore.Score' Error:Field validation for 'Score' failed on the 'required' tag"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Failure: score is -1",
			giveBody:   `{"name": "test", "score": -1}`,
			wantBody:   `{"error":"Key: 'RequestPostScore.Score' Error:Field validation for 'Score' failed on the 'min' tag"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Failure: score is 501",
			giveBody:   `{"name": "test", "score": 501}`,
			wantBody:   `{"error":"Key: 'RequestPostScore.Score' Error:Field validation for 'Score' failed on the 'max' tag"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Success: give new name and score that is valid",
			giveBody:   `{"name": "test", "score": 3}`,
			wantBody:   `{"score":3,"name":"test","rank":4}`,
			wantStatus: http.StatusOK,
		},
		{
			name:       "Success: existing name and score that is valid",
			giveBody:   `{"name": "test", "score": 33}`,
			wantBody:   `{"score":33,"name":"test","rank":1}`,
			wantStatus: http.StatusOK,
		},
		{
			name:       "Failure: existing name and score that is smaller than the previous one",
			giveBody:   `{"name": "test", "score": 22}`,
			wantBody:   `{"error":"given score 22 is smaller or equal than the already existing score 33"}`,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "Failure: existing name and score that is equal than the previous one",
			giveBody:   `{"name": "test", "score": 33}`,
			wantBody:   `{"error":"given score 33 is smaller or equal than the already existing score 33"}`,
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "", bytes.NewReader([]byte(tt.giveBody)))
			require.Nilf(t, err, "failed to create a request")
			ctx, rr := ctxAndRecorder(t, req)
			withError(controllerPostScore(service.NewService(&service.Options{
				Repo: repositories.NewRepository(db),
			}).ScoreService()))(ctx)

			assert.Equal(t, tt.wantStatus, rr.Code)
			assert.Equal(t, tt.wantBody, rr.Body.String())
		})
	}
}
