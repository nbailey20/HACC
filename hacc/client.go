package main

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

type ssmClient struct {
	Ssm *ssm.Client
}

func NewSsmClient() *ssmClient {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load AWS config, %v", err)
	}
	return &ssmClient{Ssm: ssm.NewFromConfig(cfg)}
}

func (c *ssmClient) GetParameter(name string) (string, error) {
	paramOutput, err := c.Ssm.GetParameter(context.TODO(), &ssm.GetParameterInput{
		Name:           aws.String(name),
		WithDecryption: aws.Bool(true), // Use aws.Bool to create a *bool
	})
	if err != nil {
		log.Fatalf("failed to get parameter, %v", err)
	}
	return *paramOutput.Parameter.Value, nil
}

func (c *ssmClient) PutParameter(name string, value string) error {
	_, err := c.Ssm.PutParameter(context.TODO(), &ssm.PutParameterInput{
		Name:      aws.String(name),
		Value:     aws.String(value),
		Type:      "SecureString",
		Overwrite: aws.Bool(true),
	})
	if err != nil {
		log.Fatalf("failed to put parameter, %v", err)
	}
	return nil
}

func (c *ssmClient) DeleteParameter(name string) error {
	_, err := c.Ssm.DeleteParameter(context.TODO(), &ssm.DeleteParameterInput{
		Name: aws.String(name),
	})
	if err != nil {
		var notFoundErr *types.ParameterNotFound
		if errors.As(err, &notFoundErr) {
			// Return no error if the parameter doesn't exist
			return nil
		}
		log.Fatalf("failed to delete parameter, %v", err)
	}
	return nil
}

func (c *ssmClient) GetAllParametersByPath(path string) (map[string]string, error) {
	params := make(map[string]string)
	paginator := ssm.NewGetParametersByPathPaginator(c.Ssm, &ssm.GetParametersByPathInput{
		Path: aws.String(path),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			log.Fatalf("failed to get parameters by path, %v", err)
		}
		for _, param := range page.Parameters {
			params[*param.Name] = *param.Value
		}
	}
	return params, nil
}
