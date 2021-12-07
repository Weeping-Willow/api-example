package service

import "errors"

var (
	ErrUnathorzedToken = errors.New("unathorzed token")
	ErrInvalidToken    = errors.New("invalid token")
)
