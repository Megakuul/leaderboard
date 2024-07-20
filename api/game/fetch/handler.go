package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/megakuul/leaderboard/api/game/fetch/query"
)

type FetchResponse struct {
	Message string             `json:"message"`
	Games   []query.GameOutput `json:"games"`
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
	gameid, ok := request.QueryStringParameters["gameid"]
	if ok && gameid != "" {
		games, err := query.FetchByGameId(dynamoClient, ctx, GAMETABLE, gameid)
		if err != nil {
			return nil, http.StatusNotFound, fmt.Errorf("failed to fetch data by gameid: %v", err)
		}
		return &FetchResponse{
			Message: "successfully fetched data by gameid",
			Games:   games,
		}, http.StatusOK, nil
	}
	date, ok := request.QueryStringParameters["date"]
	if ok && date != "" {
		games, err := query.FetchByDate(dynamoClient, ctx, GAMETABLE, date)
		if err != nil {
			return nil, http.StatusNotFound, fmt.Errorf("failed to fetch data by date: %v", err)
		}
		return &FetchResponse{
			Message: "successfully fetched data by date",
			Games:   games,
		}, http.StatusOK, nil
	}
	return nil, http.StatusBadRequest, fmt.Errorf("failed to fetch data: no search param provided")
}
