package router

import (
	"net/http"

	"github.com/Weeping-Willow/api-example/models"
	"github.com/Weeping-Willow/api-example/service"
	"github.com/gin-gonic/gin"
)

func controllerPostScore(score service.ScoreService) HandlerE {
	return func(c *gin.Context) error {
		var req *models.RequestPostScore
		if err := c.ShouldBindJSON(&req); err != nil {
			return newErrorBadRequest(err)
		}

		score, err := score.PostScore(req)
		if err != nil {
			return err
		}

		c.JSON(http.StatusOK, score)
		return nil
	}
}

func controllerGetScores(score service.ScoreService) HandlerE {
	return func(c *gin.Context) error {
		var req models.RequestGetScores
		if err := c.ShouldBindQuery(&req); err != nil {
			return newErrorBadRequest(err)
		}

		score, err := score.GetScores(&req)
		if err != nil {
			return err
		}

		c.JSON(http.StatusOK, score)
		return nil
	}
}
