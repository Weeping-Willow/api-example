package router

import (
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

		c.JSON(200, score)
		return nil
	}
}
