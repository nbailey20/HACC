package vault

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

type ssmClient struct {
	ssm *ssm.Client
}

func NewSsmClient(profile string) (*ssmClient, error) {
	opts := []func(*config.LoadOptions) error{
		config.WithRetryer(func() aws.Retryer {
			return retry.NewAdaptiveMode()
		}),
		config.WithRetryMaxAttempts(20),
	}
	if profile != "" {
		opts = append(
			opts,
			config.WithSharedConfigProfile(profile),
		)
	}
	cfg, err := config.LoadDefaultConfig(context.TODO(), opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %v", err)
	}
	return &ssmClient{ssm: ssm.NewFromConfig(cfg)}, nil
}

func (c *ssmClient) GetParameter(name string) (string, error) {
	paramOutput, err := c.ssm.GetParameter(context.TODO(), &ssm.GetParameterInput{
		Name:           aws.String(name),
		WithDecryption: aws.Bool(true), // Use aws.Bool to create a *bool
	})
	if err != nil {
		return "", fmt.Errorf("failed to get parameter, %v", err)
	}
	return *paramOutput.Parameter.Value, nil
}

func (c *ssmClient) PutParameter(name string, value string, key_id string) error {
	input := ssm.PutParameterInput{
		Name:      aws.String(name),
		Value:     aws.String(value),
		Type:      "SecureString",
		Overwrite: aws.Bool(true),
	}
	if key_id != "" {
		input.KeyId = aws.String(key_id)
	}
	_, err := c.ssm.PutParameter(context.TODO(), &input)
	if err != nil {
		return fmt.Errorf("failed to put parameter, %v", err)
	}
	return nil
}

func (c *ssmClient) DeleteParameter(name string) error {
	_, err := c.ssm.DeleteParameter(context.TODO(), &ssm.DeleteParameterInput{
		Name: aws.String(name),
	})
	if err != nil {
		var notFoundErr *types.ParameterNotFound
		if errors.As(err, &notFoundErr) {
			// Return no error if the parameter doesn't exist
			return nil
		}
		return fmt.Errorf("failed to delete parameter, %v", err)
	}
	return nil
}

func (c *ssmClient) GetAllParametersByPath(path string) (map[string]string, error) {
	params := make(map[string]string)
	paginator := ssm.NewGetParametersByPathPaginator(c.ssm, &ssm.GetParametersByPathInput{
		Path:      aws.String(path),
		Recursive: aws.Bool(true),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return params, fmt.Errorf("failed to get parameters by path for %s, %v", path, err)
		}
		for _, param := range page.Parameters {
			params[*param.Name] = *param.Value
		}
	}
	return params, nil
}
