package service

import "errors"

var (
	ErrUnathorzedToken           = errors.New("unathorzed token")
	ErrInvalidToken              = errors.New("invalid token")
	ErrScoreIsSmaller            = "given score %d is smaller or equal than the already existing score %d"
	ErrUpdateFailed              = errors.New("updated failed")
	ErrRankingNotFound           = errors.New("rankings not found for a score")
	ErrPaginationInvalid         = errors.New("pagination page or its size can't be smaller than 0")
	ErrPaginationPageMoreThanMax = errors.New("page can't be greater than max pages available")
)
