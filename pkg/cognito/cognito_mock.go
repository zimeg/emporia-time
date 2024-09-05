package cognito

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type CognitoMock struct {
	mock.Mock
}

func (cm *CognitoMock) GenerateTokens(ctx context.Context, username string, password string) (CognitoResponse, error) {
	args := cm.Called(ctx, username, password)
	return args.Get(0).(CognitoResponse), args.Error(1)
}

func (cm *CognitoMock) RefreshTokens(ctx context.Context, refreshToken string) (CognitoResponse, error) {
	args := cm.Called(ctx, refreshToken)
	return args.Get(0).(CognitoResponse), args.Error(1)
}
