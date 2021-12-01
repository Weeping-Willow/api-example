package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		want     *Config
		wantErr  string
	}{
		{
			name:     "Failure: no env file available, returns error",
			fileName: "",
			want:     nil,
			wantErr:  "open : no such file or directory",
		},
		{
			name:     "Success: empty env file is passed returns config with default values",
			fileName: ".env_empty",
			want: &Config{
				Port:        "8080",
				DatabaseUrl: "",
			},
			wantErr: "",
		},
		{
			name:     "Success: env file with values in it",
			fileName: ".env_filled",
			want: &Config{
				Port:        "7000",
				DatabaseUrl: "test",
			},
			wantErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.fileName)

			if tt.wantErr != "" {
				assert.Equal(t, tt.wantErr, err.Error())
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
