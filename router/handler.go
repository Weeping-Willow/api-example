package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
This wraps gin.HandlerFunc
if a controller returns error and it doesn't have implementaition of ErrorHandler,
it will return 500 error with unexpected error message
*/
type HandlerE func(*gin.Context) error

type ErrorHandler interface {
	handleError(*gin.Context) bool
}

func withError(h HandlerE) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := h(c)
		if err == nil {
			return
		}
		if er, ok := err.(ErrorHandler); ok {
			if er.handleError(c) {
				return
			}
		}
		c.JSON(http.StatusInternalServerError, errorRespose{"unexpected error"})
	}
}

type errorRespose struct {
	Err string `json:"error"`
}

func (e *errorRespose) handleError(ctx *gin.Context) bool {
	ctx.JSON(
		http.StatusUnauthorized,
		newErrorNotAuthorized(),
	)
	return true
}

func newErrorNotAuthorized() *errorRespose {
	return &errorRespose{Err: "not authorized"}
}

func (e *errorRespose) Error() string {
	return e.Err
}
