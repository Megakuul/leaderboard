package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
)

var (
	REGION                = os.Getenv("REGION")
	USERTABLE             = os.Getenv("USERTABLE")
	GAMETABLE             = os.Getenv("GAMETABLE")
	MAILTEMPLATE          = os.Getenv("MAILTEMPLATE")
	MAILSENDER            = os.Getenv("MAILSENDER")
	CONFIRM_SECRET_LENGTH = 20 // default 20
	HOURS_UNTIL_EXPIRED   = 24 // default 24
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
	sesClient := sesv2.NewFromConfig(awsConfig)

	if secretLength, err := strconv.Atoi(os.Getenv("CONFIRM_SECRET_LENGTH")); err == nil {
		CONFIRM_SECRET_LENGTH = secretLength
	}
	if hoursUntilExpired, err := strconv.Atoi(os.Getenv("HOURS_UNTIL_EXPIRED")); err == nil {
		HOURS_UNTIL_EXPIRED = hoursUntilExpired
	}

	lambda.Start(AddHandler(dynamoClient, sesClient))
	return nil
}
