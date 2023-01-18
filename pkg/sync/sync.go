package sync

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go/aws"
)

func syncSecret(
	ctx context.Context,
	secretName string,
	region string,
	secretData map[string]string) error {

	config, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return err
	}
	svc := secretsmanager.NewFromConfig(config)

	// check if secret exists
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}
	// check if secret exists
	result, err := svc.GetSecretValue(ctx, input)
	if err == nil {
		return fmt.Errorf("Secret exists. Aborting secret sync")
	}
	return nil
}
