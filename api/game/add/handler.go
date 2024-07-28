package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/megakuul/leaderboard/api/game/add/put"
	"github.com/megakuul/leaderboard/api/game/add/query"
	"github.com/megakuul/leaderboard/api/game/add/rating"
	"github.com/megakuul/leaderboard/api/game/add/sender"
)

type Participant struct {
	Username  string `json:"username"`
	Team      int    `json:"team"`
	Points    int    `json:"points"`
	Placement int    `json:"placement"`
}

type AddRequest struct {
	PlacementPoints int           `json:"placement_points"`
	Participants    []Participant `dynamodbav:"participants" json:"participants"`
}

type AddResponse struct {
	Message string `json:"message"`
	GameId  string `json:"gameid"`
}

func AddHandler(dynamoClient *dynamodb.Client, sesClient *sesv2.Client) func(context.Context, events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return func(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		response, code, err := runAddHandler(dynamoClient, sesClient, &request, ctx)
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

func runAddHandler(dynamoClient *dynamodb.Client, sesClient *sesv2.Client, request *events.APIGatewayV2HTTPRequest, ctx context.Context) (*AddResponse, int, error) {
	var req AddRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("failed to deserialize request: invalid body")
	}

	if len(req.Participants) < 2 {
		return nil, http.StatusBadRequest, fmt.Errorf("minimum number of participants is 2")
	}

	if len(req.Participants) > MAXIMUM_PARTICIPANTS {
		return nil, http.StatusBadRequest, fmt.Errorf("maximum number of participants is %d", MAXIMUM_PARTICIPANTS)
	}

	ratingInputParticipants := []rating.ParticipantInput{}
	for _, part := range req.Participants {
		user, err := query.FetchByUsername(dynamoClient, ctx, USERTABLE, part.Username)
		if err != nil {
			return nil, http.StatusNotFound, fmt.Errorf("failed to lookup %s: %v", part.Username, err)
		}
		if user.Disabled {
			return nil, http.StatusNotFound, fmt.Errorf("failed to lookup %s: user is disabled", part.Username)
		}
		ratingInputParticipants = append(ratingInputParticipants, rating.ParticipantInput{
			UserRef:   user,
			Team:      part.Team,
			Rating:    user.Elo,
			Points:    part.Points,
			Placement: part.Placement,
		})
	}

	ratingOutputParticipants := rating.CalculateRatingUpdate(ratingInputParticipants, req.PlacementPoints, MAX_LOSS_NUMBER)

	gameInputParticipants := map[string]put.ParticipantInput{}
	emailConfirmRequests := []sender.EmailConfirmRequest{}

	for _, part := range ratingOutputParticipants {
		secret := make([]byte, CONFIRM_SECRET_LENGTH)
		_, err := rand.Read(secret)
		if err != nil {
			return nil, http.StatusInternalServerError, fmt.Errorf("failed to generate confirmation secret")
		}
		base64Secret := base64.RawURLEncoding.EncodeToString(secret)

		emailConfirmRequests = append(emailConfirmRequests, sender.EmailConfirmRequest{
			Username:  part.UserRef.Username,
			Email:     part.UserRef.Email,
			Secret:    base64Secret,
			Placement: part.Placement,
			Points:    part.Points,
			EloUpdate: part.RatingUpdate,
		})

		if _, ok := gameInputParticipants[part.UserRef.Username]; ok {
			return nil, http.StatusBadRequest, fmt.Errorf("participant: %s found twice", part.UserRef.Username)
		}

		gameInputParticipants[part.UserRef.Username] = put.ParticipantInput{
			Subject:       part.UserRef.Subject,
			Username:      part.UserRef.Username,
			Underdog:      part.Underdog,
			Team:          part.Team,
			Placement:     part.Placement,
			Points:        part.Points,
			Elo:           part.Rating,
			EloUpdate:     part.RatingUpdate,
			Confirmed:     false,
			ConfirmSecret: base64Secret,
		}
	}

	expirationTime := time.Now().Add(time.Duration(HOURS_UNTIL_EXPIRED) * time.Hour)
	gameid, err := put.InsertGame(dynamoClient, ctx, GAMETABLE, gameInputParticipants, int(expirationTime.Unix()))
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("failed to insert game: %v", err)
	}

	if err := sender.SendConfirmMails(sesClient, ctx, MAILSENDER, MAILTEMPLATE, gameid, emailConfirmRequests); err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("failed to send at least one confirmation mail: %v", err)
	}

	return &AddResponse{
		Message: "successfully added game. ensure all players confirm the game to validate it",
		GameId:  gameid,
	}, http.StatusOK, nil
}
