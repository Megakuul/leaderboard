// contains wrappers for databsae update functions.
// main purpose is to abstract some boilerplate code
// away from the handler.
package update

import (
	"context"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func UpdateGame(dynamoClient *dynamodb.Client, ctx context.Context, tableName, gameid, username string, setReadonly bool) error {
	expressionAttributeNames := map[string]string{
		"#participants": "participants",
		"#username":     username,
		"#confirmed":    "confirmed",
	}
	expressionAttributeValues := map[string]types.AttributeValue{
		":confirmed": &types.AttributeValueMemberBOOL{Value: true},
	}
	var updateExpression string
	if setReadonly {
		expressionAttributeNames["#readonly"] = "readonly"
		expressionAttributeNames["#expires_in"] = "expires_in"
		expressionAttributeValues[":readonly"] = &types.AttributeValueMemberBOOL{Value: true}
		updateExpression = "SET #readonly = :readonly, #participants.#username.#confirmed = :confirmed  REMOVE #expires_in"
	} else {
		updateExpression = "SET #participants.#username.#confirmed = :confirmed"
	}

	_, err := dynamoClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"gameid": &types.AttributeValueMemberS{Value: gameid},
		},
		ConditionExpression:       aws.String("attribute_exists(gameid)"), // prevent it to upsert if not existent
		ExpressionAttributeNames:  expressionAttributeNames,
		ExpressionAttributeValues: expressionAttributeValues,
		UpdateExpression:          aws.String(updateExpression),
		ReturnValues:              types.ReturnValueAllNew,
	})
	if err != nil {
		return err
	}

	return nil
}

func UpdateUser(dynamoClient *dynamodb.Client, ctx context.Context, tableName, subject string, eloUpdate int) error {
	_, err := dynamoClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"subject": &types.AttributeValueMemberS{Value: subject},
		},
		ExpressionAttributeNames: map[string]string{
			"#elo": "elo",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":elo_update": &types.AttributeValueMemberN{Value: strconv.Itoa(eloUpdate)},
		},
		UpdateExpression: aws.String("ADD #elo :elo_update"),
		ReturnValues:     types.ReturnValueNone,
	})
	if err != nil {
		return err
	}
	return nil
}
