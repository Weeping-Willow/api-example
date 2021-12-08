package service

import "errors"

var (
	ErrUnathorzedToken = errors.New("unathorzed token")
	ErrInvalidToken    = errors.New("invalid token")
	ErrScoreIsSmaller  = "given score %d is smaller than already existing score %d"
	ErrUpdateFailed    = errors.New("updated failed")
	ErrRankingNotFound = errors.New("rankings not found for a score")
)
