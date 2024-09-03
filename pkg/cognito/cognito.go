package cognito

import (
	"context"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	awscognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

// CognitoResponse holds the authentication information from Cognito
type CognitoResponse struct {
	IdToken      *string
	RefreshToken *string
	ExpiresIn    int32
}

// Congintoir suggests expected interactions around authentication
type Congintoir interface {
	GenerateTokens(ctx context.Context, username string, password string) (CognitoResponse, error)
	RefreshTokens(ctx context.Context, refreshToken string) (CognitoResponse, error)
}

// Cognito implements the interactions for authentication
type Cognito struct {
	clientID string
	client   *awscognito.Client
}

// Create a cognito client with customized configurations
func NewClient(ctx context.Context, clientID string, region string) (*Cognito, error) {
	config, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(region))
	if err != nil {
		return &Cognito{}, err
	}
	client := awscognito.NewFromConfig(config)
	return &Cognito{
		clientID: clientID,
		client:   client,
	}, nil
}

// GenerateTokens creates new auth tokens from credentials
func (cog *Cognito) GenerateTokens(ctx context.Context, username string, password string) (
	CognitoResponse,
	error,
) {
	auth := awscognito.InitiateAuthInput{
		AuthFlow: "USER_PASSWORD_AUTH",
		AuthParameters: map[string]string{
			"USERNAME": username,
			"PASSWORD": password,
		},
		ClientId: &cog.clientID,
	}
	user, err := cog.client.InitiateAuth(ctx, &auth)
	if err != nil {
		return CognitoResponse{}, err
	}
	return CognitoResponse{
		IdToken:      user.AuthenticationResult.IdToken,
		RefreshToken: user.AuthenticationResult.RefreshToken,
		ExpiresIn:    user.AuthenticationResult.ExpiresIn,
	}, nil
}

// RefreshTokens regenerates auth tokens from the refresh token
func (cog *Cognito) RefreshTokens(ctx context.Context, refreshToken string) (
	CognitoResponse,
	error,
) {
	auth := awscognito.InitiateAuthInput{
		AuthFlow: "REFRESH_TOKEN_AUTH",
		AuthParameters: map[string]string{
			"REFRESH_TOKEN": refreshToken,
		},
		ClientId: &cog.clientID,
	}
	user, err := cog.client.InitiateAuth(ctx, &auth)
	if err != nil {
		return CognitoResponse{}, err
	}
	return CognitoResponse{
		IdToken:      user.AuthenticationResult.IdToken,
		RefreshToken: user.AuthenticationResult.RefreshToken,
		ExpiresIn:    user.AuthenticationResult.ExpiresIn,
	}, nil
}
