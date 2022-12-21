package main

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

// EmporiaClient is the AWS Cognito client id used by Emporia
var EmporiaClient = "4qte47jbstod8apnfic0bunmrq"

// GenerateTokens creates new auth tokens from credentials
func GenerateTokens(username string, password string) *cognito.AuthenticationResultType {
	client := createCognitoClient()

	auth := &cognito.InitiateAuthInput{
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(username),
			"PASSWORD": aws.String(password),
		},
		ClientId: aws.String(EmporiaClient),
	}

	user, err := client.InitiateAuth(auth)
	if err != nil {
		log.Fatalf("Failed to authenticate with Cognito: %s\n", err)
	}

	return user.AuthenticationResult
}

// RefreshTokens regenerates auth tokens from the refresh token
func RefreshTokens(token string) *cognito.AuthenticationResultType {
	client := createCognitoClient()

	auth := &cognito.InitiateAuthInput{
		AuthFlow: aws.String("REFRESH_TOKEN_AUTH"),
		AuthParameters: map[string]*string{
			"REFRESH_TOKEN": aws.String(token),
		},
		ClientId: aws.String(EmporiaClient),
	}

	user, err := client.InitiateAuth(auth)
	if err != nil {
		log.Fatalf("Failed to re-authenticate with Cognito: %s\n", err)
	}

	return user.AuthenticationResult
}

// SaveTokens stores newly gathered auth tokens in the config
func (conf *EmporiaConfig) SaveTokens(auth *cognito.AuthenticationResultType) {
	conf.EmporiaToken = *auth.IdToken

	if auth.RefreshToken != nil {
		conf.EmporiaRefresh = *auth.RefreshToken
	}

	lifespan := time.Duration(*auth.ExpiresIn)
	conf.EmporiaExpires = time.Now().Add(time.Second * lifespan).UTC()
}

// createCognitoClient creates a configured identity provider
func createCognitoClient() *cognito.CognitoIdentityProvider {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-2"),
	})

	if err != nil {
		log.Panicf("Failed to create an AWS session: %s\n", err)
	}

	client := cognito.New(sess)
	return client
}
