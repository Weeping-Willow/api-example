package router

import (
	"net/http"
	"testing"

	"github.com/Weeping-Willow/api-example/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_integration_middlewareTokenAuth(t *testing.T) {
	tests := []struct {
		name                    string
		skipAuthorizationHeader bool
		authorization           string
		wantBody                string
		wantStatus              int
	}{
		{
			name:                    "Negative: No Authorization header",
			skipAuthorizationHeader: true,
			wantBody:                `{"error":"not authorized"}`,
			wantStatus:              http.StatusUnauthorized,
		},
		{
			name:          "Negative: Invalid Authorization header",
			wantBody:      `{"error":"not authorized"}`,
			wantStatus:    http.StatusUnauthorized,
			authorization: "Bearer",
		},
		{
			name:          "Negative: Invalid Authorization token",
			wantBody:      `{"error":"not authorized"}`,
			wantStatus:    http.StatusUnauthorized,
			authorization: "Bearer invalid",
		},
		{
			name:          "Success: Correct Authorization token",
			wantBody:      ``,
			wantStatus:    http.StatusOK,
			authorization: "Bearer complicated-token",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "", nil)
			require.Nilf(t, err, "failed to create a request")
			if !tt.skipAuthorizationHeader {
				req.Header.Add("Authorization", tt.authorization)
			}
			ctx, rr := ctxAndRecorder(t, req)
			withError(middlewareTokenAuth(service.NewService(&service.Options{}).TokenService()))(ctx)

			assert.Equal(t, tt.wantStatus, rr.Code)
			assert.Equal(t, tt.wantBody, rr.Body.String())
		})
	}
}
