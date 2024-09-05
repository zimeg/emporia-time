package config

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zimeg/emporia-time/pkg/api"
	"github.com/zimeg/emporia-time/pkg/cognito"
)

func TestLoad(t *testing.T) {
	mockIDToken := "eyJ-example-token"
	mockRefreshToken := "eyJ-example-refresh"

	tests := map[string]struct {
		mockConfigFile                 string
		mockFlags                      Flags
		mockGenerateTokensResponse     cognito.CognitoResponse
		mockGenerateTokensError        error
		mockGetCustomerDevicesResponse []api.Device
		mockGetCustomerDevicesError    error
		mockRefreshTokensResponse      cognito.CognitoResponse
		mockRefreshTokensError         error
		expectedConfig                 Config
		expectedError                  error
	}{
		"loads the saved and valid credentials into configurations": {
			mockConfigFile: `{
                "Device": "123456",
                "Tokens": {
                    "IdToken": "eyJ-example-001",
                    "RefreshToken": "eyJ-example-002",
                    "ExpiresAt": "2222-02-22T22:22:22Z"
                }
            }`,
			mockGetCustomerDevicesResponse: []api.Device{
				{
					DeviceGid: 123456,
				},
			},
			expectedConfig: Config{
				Device: "123456",
				Tokens: TokenSet{
					IdToken:      "eyJ-example-001",
					RefreshToken: "eyJ-example-002",
					ExpiresAt:    time.Date(2222, 2, 22, 22, 22, 22, 0, time.UTC),
				},
			},
		},
		"writes configured authentication from provided credentials": {
			mockFlags: Flags{
				Username: "watt@example.com",
				Password: "joules123",
			},
			mockGenerateTokensResponse: cognito.CognitoResponse{
				IdToken:      &mockIDToken,
				RefreshToken: &mockRefreshToken,
				ExpiresIn:    1,
			},
			mockGetCustomerDevicesResponse: []api.Device{
				{
					DeviceGid: 1000001,
				},
			},
			expectedConfig: Config{
				Device: "1000001",
				Tokens: TokenSet{
					IdToken:      mockIDToken,
					RefreshToken: mockRefreshToken,
				},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			fs := afero.NewMemMapFs()
			cog := &cognito.CognitoMock{}
			cog.On("GenerateTokens", mock.Anything, mock.Anything, mock.Anything).
				Return(tt.mockGenerateTokensResponse, tt.mockGenerateTokensError)
			cog.On("RefreshTokens", mock.Anything, mock.Anything).
				Return(tt.mockRefreshTokensResponse, tt.mockRefreshTokensError)
			req := &api.EmporiaMock{}
			req.On("GetCustomerDevices").
				Return(tt.mockGetCustomerDevicesResponse, tt.mockGetCustomerDevicesError)
			req.On("SetToken", mock.Anything)
			req.On("SetDevice", mock.Anything)
			dir, err := os.UserHomeDir()
			require.NoError(t, err)
			if tt.mockConfigFile != "" {
				settings, err := fs.Create(dir + "/.config/etime/settings.json")
				require.NoError(t, err)
				_, err = settings.WriteString(tt.mockConfigFile)
				require.NoError(t, err)
			}
			os.Setenv("EMPORIA_USERNAME", tt.mockFlags.Username) // FIXME: use flags!
			os.Setenv("EMPORIA_PASSWORD", tt.mockFlags.Password) // FIXME: use flags!
			cfg, err := Load(ctx, cog, fs, req, tt.mockFlags)
			if tt.expectedError != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedConfig.Device, cfg.Device)
				assert.Equal(t, tt.expectedConfig.Tokens.IdToken, cfg.Tokens.IdToken)
				assert.Equal(t, tt.expectedConfig.Tokens.RefreshToken, cfg.Tokens.RefreshToken)
				assert.Greater(t, cfg.Tokens.ExpiresAt, time.Now())
				assert.Equal(t, dir+"/.config/etime/settings.json", cfg.path)
				req.AssertCalled(t, "SetDevice", tt.expectedConfig.Device)
				req.AssertCalled(t, "SetToken", tt.expectedConfig.Tokens.IdToken)
				actualConfigFile, err := afero.ReadFile(fs, cfg.path)
				require.NoError(t, err)
				expectedConfigFile, err := json.MarshalIndent(cfg, "", "\t")
				require.NoError(t, err)
				assert.Equal(t, expectedConfigFile, actualConfigFile)
			}
		})
	}
}