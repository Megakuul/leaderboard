package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var (
	REGION    = os.Getenv("AWS_REGION")
	USERTABLE = os.Getenv("USERTABLE")
	GAMETABLE = os.Getenv("GAMETABLE")
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("ERROR INITIALIZATION: %v\n", err)
	}
}

func run() error {
	awsConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(REGION))
	if err != nil {
		return fmt.Errorf("failed to load aws config: %v", err)
	}
	dynamoClient := dynamodb.NewFromConfig(awsConfig)

	lambda.Start(ConfirmHandler(dynamoClient))
	return nil
}
