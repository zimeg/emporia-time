package emporia

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

// EmporiaCognitoClientID is the AWS Cognito client ID used by Emporia
const EmporiaCognitoClientID = "4qte47jbstod8apnfic0bunmrq"

// EmporiaCognitoResponse holds the authentication information from Cognito
type EmporiaCognitoResponse struct {
	IdToken      *string
	RefreshToken *string
	ExpiresIn    *int64
}

// GenerateTokens creates new auth tokens from credentials
func GenerateTokens(credentials EmporiaCredentials) (EmporiaCognitoResponse, error) {
	auth := cognito.InitiateAuthInput{
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(credentials.Username),
			"PASSWORD": aws.String(credentials.Password),
		},
		ClientId: aws.String(EmporiaCognitoClientID),
	}

	if client, err := createCognitoClient(); err != nil {
		return EmporiaCognitoResponse{}, err
	} else if user, err := client.InitiateAuth(&auth); err != nil {
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
	auth := cognito.InitiateAuthInput{
		AuthFlow: aws.String("REFRESH_TOKEN_AUTH"),
		AuthParameters: map[string]*string{
			"REFRESH_TOKEN": aws.String(refreshToken),
		},
		ClientId: aws.String(EmporiaCognitoClientID),
	}

	if client, err := createCognitoClient(); err != nil {
		return EmporiaCognitoResponse{}, err
	} else if user, err := client.InitiateAuth(&auth); err != nil {
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
func createCognitoClient() (*cognito.CognitoIdentityProvider, error) {
	cfg := aws.Config{
		Region: aws.String("us-east-2"),
	}
	if sess, err := session.NewSession(&cfg); err != nil {
		return &cognito.CognitoIdentityProvider{}, err
	} else {
		return cognito.New(sess), nil
	}
}
