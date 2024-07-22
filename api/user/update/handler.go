package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/megakuul/leaderboard/api/user/update/update"
)

type UpdateResponse struct {
	Message     string            `json:"message"`
	UpdatedUser update.UserOutput `json:"updated_user"`
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

	sub := request.RequestContext.Authorizer.JWT.Claims["sub"]
	if sub == "" {
		return nil, http.StatusUnprocessableEntity, fmt.Errorf("invalid sub claim in the ID token")
	}

	user, err := update.UpsertUser(dynamoClient, ctx, BASEELO, USERTABLE, sub, REGION, request.RequestContext.Authorizer.JWT.Claims)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("failed to upsert user: %v", err)
	}

	return &UpdateResponse{
		Message:     "successfully updated user",
		UpdatedUser: *user,
	}, http.StatusOK, nil
}
