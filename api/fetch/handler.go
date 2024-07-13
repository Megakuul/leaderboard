package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/megakuul/leaderboard/api/fetch/query"
)

type FetchResponse struct {
	Message    string       `json:"message"`
	NewPageKey string       `json:"newpagekey"`
	Users      []query.User `json:"users"`
}

func FetchHandler(dynamoClient *dynamodb.Client) func(context.Context, events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return func(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		response, code, err := runFetchHandler(dynamoClient, &request, ctx)
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

func runFetchHandler(dynamoClient *dynamodb.Client, request *events.APIGatewayV2HTTPRequest, ctx context.Context) (*FetchResponse, int, error) {
	username, ok := request.QueryStringParameters["username"]
	if ok && username != "" {
		users, err := query.FetchByUsername(dynamoClient, ctx, USERTABLE, username)
		if err != nil {
			return nil, http.StatusNotFound, fmt.Errorf("failed to fetch data by username: %v", err)
		}
		return &FetchResponse{
			Message: "successfully fetched data by username",
			Users:   users,
		}, http.StatusOK, nil
	}
	elo, ok := request.QueryStringParameters["elo"]
	if ok && elo != "" {
		users, err := query.FetchByElo(dynamoClient, ctx, USERTABLE, elo)
		if err != nil {
			return nil, http.StatusNotFound, fmt.Errorf("failed to fetch data by elo: %v", err)
		}
		return &FetchResponse{
			Message: "successfully fetched data by elo",
			Users:   users,
		}, http.StatusOK, nil
	}
	lastPageKey, ok := request.QueryStringParameters["lastpagekey"]
	if !ok {
		lastPageKey = ""
	}
	users, newPageKey, err := query.FetchByPage(dynamoClient, ctx, USERTABLE, lastPageKey)
	if err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("failed to fetch data by elo: %v", err)
	}
	return &FetchResponse{
		Message:    "successfully fetched data by elo",
		NewPageKey: newPageKey,
		Users:      users,
	}, http.StatusOK, nil
}
