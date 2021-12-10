package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
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
		log.Err(err)
		c.JSON(http.StatusInternalServerError, errorRespose{Err: err.Error()})
	}
}

type errorRespose struct {
	Err string `json:"error"`
}

func (e *errorRespose) handleError(ctx *gin.Context) bool {
	ctx.JSON(
		http.StatusUnauthorized,
		newErrorNotAuthorized(ctx),
	)
	return true
}

func newErrorNotAuthorized(c *gin.Context) *errorRespose {
	c.Abort()
	return &errorRespose{Err: "not authorized"}
}

func (e *errorRespose) Error() string {
	return e.Err
}

type errorBadRequest struct {
	Err string `json:"error"`
}

func (e *errorBadRequest) handleError(ctx *gin.Context) bool {
	ctx.JSON(
		http.StatusBadRequest,
		newErrorBadRequest(e),
	)
	return true
}

func newErrorBadRequest(err error) *errorBadRequest {
	return &errorBadRequest{Err: err.Error()}
}

func (e *errorBadRequest) Error() string {
	return e.Err
}
