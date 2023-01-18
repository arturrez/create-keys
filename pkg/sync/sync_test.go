package sync

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type mockASM struct {
	returnErrorOnCreate   bool
	returnErrorOnDescribe bool
}

func (m mockASM) CreateSecret(ctx context.Context, params *secretsmanager.CreateSecretInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.CreateSecretOutput, error) {
	fmt.Println("CreateSecret")
	r := &secretsmanager.CreateSecretOutput{
		Name: params.Name,
		ARN:  aws.String("123456:fake-arn"),
	}
	if m.returnErrorOnCreate {
		return r, fmt.Errorf(*params.Name)
	}
	return r, nil
}

func (m mockASM) DescribeSecret(ctx context.Context, params *secretsmanager.DescribeSecretInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.DescribeSecretOutput, error) {
	fmt.Println("DescribeSecret " + *params.SecretId)
	r := &secretsmanager.DescribeSecretOutput{}
	if m.returnErrorOnDescribe {
		return r, fmt.Errorf(*params.SecretId)
	}
	return r, nil
}

func TestCreateSecret(t *testing.T) {
	tests := map[string]struct {
		inSecretName  string
		canCreate     bool
		canDescribe   bool
		inSecretMap   map[string]string
		expectedError bool
	}{
		"should return error if secret already exists": {
			inSecretName:  "Exists",
			canCreate:     true,
			canDescribe:   false,
			inSecretMap:   map[string]string{"key": "value"},
			expectedError: true,
		},

		"should return no error if successful": {
			inSecretName:  "notExists",
			canCreate:     false,
			canDescribe:   true,
			inSecretMap:   map[string]string{"key": "value"},
			expectedError: false,
		},
	}
	var m mockASM
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m.returnErrorOnCreate = tc.canCreate
			m.returnErrorOnDescribe = tc.canDescribe
			mockSecretsManager := SecretsManager{
				client: m,
				ctx:    context.TODO(),
				region: "us-east-1"}

			_, err := mockSecretsManager.createSecret(tc.inSecretName, tc.inSecretMap)
			fmt.Println(err)
			// THEN
			require.Equal(t, tc.expectedError, err != nil)
		})
	}
}
