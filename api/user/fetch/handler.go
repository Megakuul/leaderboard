package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/megakuul/leaderboard/api/user/fetch/query"
)

type FetchResponse struct {
	Message    string             `json:"message"`
	NewPageKey string             `json:"newpagekey"`
	Users      []query.UserOutput `json:"users"`
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
	region := request.QueryStringParameters["region"]
	if region == "" {
		region = REGION
	}

	pageSizeStr := request.QueryStringParameters["pagesize"]
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = query.MAX_PAGESIZE
	}

	lastPageKey, ok := request.QueryStringParameters["lastpagekey"]
	if !ok {
		lastPageKey = ""
	}

	username := request.QueryStringParameters["username"]
	if username != "" {
		users, err := query.FetchByUsername(dynamoClient, ctx, USERTABLE, username)
		if err != nil {
			return nil, http.StatusNotFound, fmt.Errorf("failed to fetch data by username: %v", err)
		}
		return &FetchResponse{
			Message: "successfully fetched data by username",
			Users:   users,
		}, http.StatusOK, nil
	}

	elo := request.QueryStringParameters["elo"]
	if elo != "" {
		users, err := query.FetchByElo(dynamoClient, ctx, USERTABLE, int32(pageSize), region, elo)
		if err != nil {
			return nil, http.StatusNotFound, fmt.Errorf("failed to fetch data by elo: %v", err)
		}
		return &FetchResponse{
			Message: "successfully fetched data by elo",
			Users:   users,
		}, http.StatusOK, nil
	}

	users, newPageKey, err := query.FetchByPage(dynamoClient, ctx, USERTABLE, int32(pageSize), lastPageKey, region)
	if err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("failed to fetch data by page: %v", err)
	}

	return &FetchResponse{
		Message:    "successfully fetched data by page",
		NewPageKey: newPageKey,
		Users:      users,
	}, http.StatusOK, nil
}
