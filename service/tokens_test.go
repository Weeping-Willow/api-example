package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_tokenService_Check(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantErr error
	}{
		{
			name:    "Failure: token is empty",
			token:   "",
			wantErr: ErrInvalidToken,
		},
		{
			name:    "Failure: token is not valid",
			token:   "Bearer randomFakeToken",
			wantErr: ErrUnathorzedToken,
		},
		{
			name:    "Success: token is valid",
			token:   "Bearer complicated-token",
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &tokenService{
				commonServices: commonServices{},
			}

			err := tr.Check(tt.token)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
