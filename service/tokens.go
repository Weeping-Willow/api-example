package service

import "strings"

type TokenService interface {
	Check(token string) error
}

type tokenService struct {
	commonService
}

func newTokenService(opts *Options) *tokenService {
	return &tokenService{
		commonService: commonService{Repo: opts.Repo, Config: opts.Config},
	}
}

func (s *service) TokenService() TokenService {
	return s.tokensService
}

func (t *tokenService) Check(token string) error {
	if !strings.Contains(token, "Bearer") {
		return ErrInvalidToken
	}
	token = strings.Replace(token, "Bearer ", "", 1)

	if token == "complicated-token" {
		return nil
	}

	return ErrUnathorzedToken
}
