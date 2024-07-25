// contains wrappers for database put functions.
// main purpose is to abstract some boilerplate code
// away from the handler.
package put

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type ParticipantInput struct {
	Subject       string `dynamodbav:"subject"`
	Username      string `dynamodbav:"username"`
	Underdog      bool   `dynamodbav:"underdog"`
	Team          int    `dynamodbav:"team"`
	Placement     int    `dynamodbav:"placement"`
	Points        int    `dynamodbav:"points"`
	Elo           int    `dynamodbav:"elo"`
	EloUpdate     int    `dynamodbav:"elo_update"`
	Confirmed     bool   `dynamodbav:"confirmed"`
	ConfirmSecret string `dynamodbav:"confirm_secret"`
}

type GameInput struct {
	GameId      string                      `dynamodbav:"gameid"`
	Date        string                      `dynamodbav:"game_date"`
	ExpiresIn   int                         `dynamodbav:"expires_in"`
	Readonly    bool                        `dynamodbav:"readonly"`
	Partcipants map[string]ParticipantInput `dynamodbav:"participants"`
}

func InsertGame(dynamoClient *dynamodb.Client, ctx context.Context, tableName string, participants map[string]ParticipantInput, expirationTime int) (string, error) {
	gameId := uuid.New().String()

	gameInput := GameInput{
		GameId:      gameId,
		Date:        time.Now().Format("2006-01-02"),
		Readonly:    false,
		ExpiresIn:   expirationTime,
		Partcipants: participants,
	}
	gameInputSerialized, err := attributevalue.MarshalMap(&gameInput)
	if err != nil {
		return "", fmt.Errorf("failed to serialize put input")
	}

	_, err = dynamoClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:    aws.String(tableName),
		Item:         gameInputSerialized,
		ReturnValues: types.ReturnValueNone,
	})
	if err != nil {
		return "", err
	}

	return gameId, nil
}
