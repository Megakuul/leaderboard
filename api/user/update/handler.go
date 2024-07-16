package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type UpdateResponse struct {
	Message string `json:"message"`
}

func UpdateHandler(dynamoClient *dynamodb.Client) func(context.Context, events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return func(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		response, code, err := runUpdatehandler(dynamoClient, &request, ctx)
		if err != nil {
			return events.APIGatewayV2HTTPResponse{
				StatusCode: code,
				Headers:    map[string]string{"Content-Type": "text/plain"},
				Body:       err.Error(),
			}, nil
		}
		serializedResponse, err := json.Marshal(&response)
		if err != nil {
			return events.APIGatewayV2HTTPResponse{
				StatusCode: http.StatusInternalServerError,
				Headers:    map[string]string{"Content-Type": "text/plain"},
				Body:       "failed to serialize response",
			}, nil
		}
		return events.APIGatewayV2HTTPResponse{
			StatusCode: code,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       string(serializedResponse),
		}, nil
	}
}

func runUpdatehandler(dynamoClient *dynamodb.Client, request *events.APIGatewayV2HTTPRequest, ctx context.Context) (*UpdateResponse, int, error) {

}
