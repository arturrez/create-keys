package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"github.com/aws/aws-sdk-go/aws"
)

const secretTag string = "avalanchego-secret-operator"

type api interface {
	CreateSecret(ctx context.Context, params *secretsmanager.CreateSecretInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.CreateSecretOutput, error)
	DescribeSecret(ctx context.Context, params *secretsmanager.DescribeSecretInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.DescribeSecretOutput, error)
}

type SecretsManager struct {
	client api
	region string
	ctx    context.Context
}

// New returns a SecretsManager configured with the default session.
func New(ctx context.Context,
	region string,
) (*SecretsManager, error) {
	config, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, err
	}
	svc := secretsmanager.NewFromConfig(config)
	return &SecretsManager{
		client: svc,
		region: region,
		ctx:    ctx,
	}, nil
}

// todo add KMS support
func (s *SecretsManager) createSecret(
	secretID string,
	secretData map[string]string) (string, error) {

	secretJson, err := json.Marshal(secretData)
	if err != nil {
		return "", err
	}
	secretDescription := "Managed by " + secretTag
	input := &secretsmanager.CreateSecretInput{
		Name:                        aws.String(secretID),
		ForceOverwriteReplicaSecret: false,
		Description:                 &secretDescription,
		SecretBinary:                secretJson,
		Tags:                        s.generateTags(),
	}
	// check if secret exists
	_, err = s.client.DescribeSecret(s.ctx, &secretsmanager.DescribeSecretInput{SecretId: aws.String(secretID)})
	if err == nil {
		return "", fmt.Errorf("secret exists. Aborting secret sync")
	}

	//write secretData
	out, err := s.client.CreateSecret(s.ctx, input)
	if err != nil {
		return "", err
	}

	return *out.ARN, nil
}

func (s *SecretsManager) generateTags() []types.Tag {
	tags := []types.Tag{}
	timestamp := time.Now().UTC().Format(time.UnixDate)
	return append(tags, types.Tag{Key: aws.String(secretTag), Value: aws.String(timestamp)})
}
