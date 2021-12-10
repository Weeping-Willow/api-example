package router

import (
	"github.com/Weeping-Willow/api-example/service"
	"github.com/gin-gonic/gin"
)

func middlewareTokenAuth(token service.TokenService) HandlerE {
	return func(c *gin.Context) error {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			return newErrorNotAuthorized(c)
		}

		if token.Check(auth) != nil {
			return newErrorNotAuthorized(c)
		}

		c.Next()
		return nil
	}
}
