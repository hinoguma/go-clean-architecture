package config

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"os"
)

func GetAwsConfig() (aws.Config, error) {
	return config.LoadDefaultConfig(context.TODO())
}
func GetAWSRegion() string {
	return os.Getenv("AWS_DEFAULT_REGION")
}

func GetCognitoAppClientID() string {
	return os.Getenv("COGNITO_APP_CLIENT_ID")
}

func GetCognitoUserPoolID() string {
	return os.Getenv("COGNITO_USER_POOL_ID")
}
