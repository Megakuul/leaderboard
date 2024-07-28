package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/megakuul/leaderboard/api/game/confirm/query"
	"github.com/megakuul/leaderboard/api/game/confirm/update"
)

func ConfirmHandler(dynamoClient *dynamodb.Client) func(context.Context, events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return func(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		response, code, err := runConfirmHandler(dynamoClient, &request, ctx)
		if err != nil {
			return events.APIGatewayV2HTTPResponse{
				StatusCode: code,
				Headers:    map[string]string{"Content-Type": "text/plain"},
				Body:       err.Error(),
			}, nil
		}
		return events.APIGatewayV2HTTPResponse{
			StatusCode: code,
			Headers:    map[string]string{"Content-Type": "text/plain"},
			Body:       response,
		}, nil
	}
}

func runConfirmHandler(dynamoClient *dynamodb.Client, request *events.APIGatewayV2HTTPRequest, ctx context.Context) (string, int, error) {
	gameid, ok := request.QueryStringParameters["gameid"]
	if !ok || gameid == "" {
		return "", http.StatusBadRequest, fmt.Errorf("missing query parameter 'gameid'")
	}

	username, ok := request.QueryStringParameters["username"]
	if !ok || username == "" {
		return "", http.StatusBadRequest, fmt.Errorf("missing query parameter 'username'")
	}

	code, ok := request.QueryStringParameters["code"]
	if !ok || code == "" {
		return "", http.StatusBadRequest, fmt.Errorf("missing query parameter 'code'")
	}

	game, err := query.FetchById(dynamoClient, ctx, GAMETABLE, gameid)
	if err != nil {
		return "", http.StatusNotFound, fmt.Errorf("failed to confirm: %v", err)
	}

	if game.Readonly {
		return "", http.StatusBadRequest, fmt.Errorf("the game was already confirmed by all participants and is now readonly")
	}

	subject := ""
	confirmCount := 0
	for _, part := range game.Participants {
		if part.Confirmed {
			confirmCount++
			continue
		} else if part.Username == username {
			if part.ConfirmSecret != code {
				return "", http.StatusForbidden, fmt.Errorf("invalid confirmation code")
			}
			subject = part.Subject
			confirmCount++
			continue
		}
	}
	if subject == "" {
		return "", http.StatusNotFound, fmt.Errorf("user not found or already confirmed in specified game")
	}

	if confirmCount != len(game.Participants) {
		if err := update.UpdateGame(dynamoClient, ctx, GAMETABLE, gameid, username, false); err != nil {
			return "", http.StatusBadRequest, fmt.Errorf("failed to update game: %v", err)
		}
		return fmt.Sprintf("successfully confirmed game %s", gameid), http.StatusOK, nil
	}

	userUpdateFailure := false
	for _, part := range game.Participants {
		err = update.UpdateUser(dynamoClient, ctx, USERTABLE, part.Subject, part.EloUpdate)
		if err != nil {
			userUpdateFailure = true
		}
	}

	if err := update.UpdateGame(dynamoClient, ctx, GAMETABLE, gameid, username, true); err != nil {
		return "", http.StatusBadRequest, fmt.Errorf("failed to update game: %v", err)
	}

	if userUpdateFailure {
		return "", http.StatusInternalServerError, fmt.Errorf(
			"game update successful, but one or more user updates failed. If points are missing, please contact an administrator")
	} else {
		return fmt.Sprintf("successfully confirmed game %s", gameid), http.StatusOK, nil
	}
}
