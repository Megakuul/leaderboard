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

type ParticipantOutput struct {
	Username  string `dynamodbav:"username" json:"username"`
	Placement int    `dynamodbav:"placement" json:"placement"`
	Points    int    `dynamodbav:"points" json:"points"`
	Elo       int    `dynamodbav:"elo" json:"elo"`
	EloUpdate int    `dynamodbav:"elo_update" json:"elo_update"`
	Confirmed bool   `dynamodbav:"confirmed" json:"confirmed"`
}

type GameOutput struct {
	GameId      string                      `dynamodbav:"gameid" json:"gameid"`
	Date        string                      `dynamodbav:"game_date" json:"date"`
	Readonly    bool                        `dynamodbav:"readonly" json:"readonly"`
	ExpiresIn   int                         `dynamodbav:"expires_in" json:"expires_in"`
	Partcipants map[string]ParticipantInput `dynamodbav:"participants" json:"participants"`
}

func InsertGame(dynamoClient *dynamodb.Client, ctx context.Context, tableName string, participants map[string]ParticipantInput, expirationTime int) (*GameOutput, error) {
	gameInput := GameInput{
		GameId:      uuid.New().String(),
		Date:        time.Now().Format("2006-01-02"),
		Readonly:    false,
		ExpiresIn:   expirationTime,
		Partcipants: participants,
	}
	gameInputSerialized, err := attributevalue.MarshalMap(&gameInput)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize put input")
	}

	output, err := dynamoClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:    aws.String(tableName),
		Item:         gameInputSerialized,
		ReturnValues: types.ReturnValueAllNew,
	})
	if err != nil {
		return nil, err
	}

	var gameOutput GameOutput
	if err := attributevalue.UnmarshalMap(output.Attributes, &gameOutput); err != nil {
		return nil, fmt.Errorf("failed to deserialize put output")
	}

	return &gameOutput, nil
}
