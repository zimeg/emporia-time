package emporia

import (
	"context"

	config "github.com/aws/aws-sdk-go-v2/config"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

// EmporiaCognitoClientID is the AWS Cognito client ID used by Emporia
var EmporiaCognitoClientID string = "4qte47jbstod8apnfic0bunmrq"

// EmporiaCognitoResponse holds the authentication information from Cognito
type EmporiaCognitoResponse struct {
	IdToken      *string
	RefreshToken *string
	ExpiresIn    int32
}

// GenerateTokens creates new auth tokens from credentials
func GenerateTokens(credentials EmporiaCredentials) (EmporiaCognitoResponse, error) {
	ctx := context.Background()
	auth := cognito.InitiateAuthInput{
		AuthFlow: "USER_PASSWORD_AUTH",
		AuthParameters: map[string]string{
			"USERNAME": credentials.Username,
			"PASSWORD": credentials.Password,
		},
		ClientId: &EmporiaCognitoClientID,
	}
	if client, err := createCognitoClient(); err != nil {
		return EmporiaCognitoResponse{}, err
	} else if user, err := client.InitiateAuth(ctx, &auth); err != nil {
		return EmporiaCognitoResponse{}, err
	} else {
		return EmporiaCognitoResponse{
			IdToken:      user.AuthenticationResult.IdToken,
			RefreshToken: user.AuthenticationResult.RefreshToken,
			ExpiresIn:    user.AuthenticationResult.ExpiresIn,
		}, nil
	}
}

// RefreshTokens regenerates auth tokens from the refresh token
func RefreshTokens(refreshToken string) (EmporiaCognitoResponse, error) {
	ctx := context.Background()
	auth := cognito.InitiateAuthInput{
		AuthFlow: "REFRESH_TOKEN_AUTH",
		AuthParameters: map[string]string{
			"REFRESH_TOKEN": refreshToken,
		},
		ClientId: &EmporiaCognitoClientID,
	}
	if client, err := createCognitoClient(); err != nil {
		return EmporiaCognitoResponse{}, err
	} else if user, err := client.InitiateAuth(ctx, &auth); err != nil {
		return EmporiaCognitoResponse{}, err
	} else {
		return EmporiaCognitoResponse{
			IdToken:      user.AuthenticationResult.IdToken,
			RefreshToken: user.AuthenticationResult.RefreshToken,
			ExpiresIn:    user.AuthenticationResult.ExpiresIn,
		}, nil
	}
}

// createCognitoClient creates a configured identity provider
func createCognitoClient() (*cognito.Client, error) {
	ctx := context.Background()
	config, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-2"))
	if err != nil {
		return &cognito.Client{}, err
	}
	return cognito.NewFromConfig(config), nil
}
