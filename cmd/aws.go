package cmd

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

var client *ecs.Client
var awsConfig aws.Config

func makeClient(region string) {
	var err error

	if region == "" {
		awsConfig, err = config.LoadDefaultConfig(context.TODO())
	} else {
		awsConfig, err = config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	}

	if err != nil {
		log.Fatalf("failed to initialize AWS client: %v", err)
	}

	client = ecs.NewFromConfig(awsConfig)
}
