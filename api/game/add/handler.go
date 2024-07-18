package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Participant struct {
	Username  string `json:"username"`
	Points    int    `json:"points"`
	Placement int    `json:"placement"`
}

type AddRequest struct {
	Participants []Participant `dynamodbav:"participants" json:"participants"`
}

type AddResponse struct {
	Message string `json:"message"`
	GameId  string `json:"gameid"`
}

func AddHandler(dynamoClient *dynamodb.Client) func(context.Context, events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return func(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		response, code, err := runAddHandler(dynamoClient, &request, ctx)
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

func runAddHandler(dynamoClient *dynamodb.Client, request *events.APIGatewayV2HTTPRequest, ctx context.Context) (*AddResponse, int, error) {
	var req AddRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("failed to deserialize request: invalid body")
	}

}
