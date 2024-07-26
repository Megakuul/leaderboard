// contains wrappers for databsae update functions.
// main purpose is to abstract some boilerplate code
// away from the handler.
package update

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type UserInput struct {
	Title    string `json:"title"`
	IconURL  string `json:"iconurl"`
	Disabled bool   `json:"disabled"`
}

type UserOutput struct {
	Username string `dynamodbav:"username" json:"username"`
	Disabled bool   `dynamodbav:"disabled" json:"disabled"`
	Region   string `dynamodbav:"user_region" json:"region"`
	Title    string `dynamodbav:"title" json:"title"`
	Email    string `dynamodbav:"email" json:"email"`
	IconURL  string `dynamodbav:"iconurl" json:"iconurl"`
	Elo      int    `dynamodbav:"elo" json:"elo"`
}

func UpsertUser(dynamoClient *dynamodb.Client, ctx context.Context, baseElo, tableName, subject, region string, claims map[string]string, userUpdate *UserInput) (*UserOutput, error) {
	output, err := dynamoClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"subject": &types.AttributeValueMemberS{Value: subject},
		},
		ExpressionAttributeNames: map[string]string{
			"#username":    "username",
			"#title":       "title",
			"#disabled":    "disabled",
			"#iconurl":     "iconurl",
			"#email":       "email",
			"#user_region": "user_region",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":username":    &types.AttributeValueMemberS{Value: claims["preferred_username"]},
			":disabled":    &types.AttributeValueMemberBOOL{Value: userUpdate.Disabled},
			":title":       &types.AttributeValueMemberS{Value: userUpdate.Title},
			":iconurl":     &types.AttributeValueMemberS{Value: userUpdate.IconURL},
			":email":       &types.AttributeValueMemberS{Value: claims["email"]},
			":user_region": &types.AttributeValueMemberS{Value: region},
		},
		UpdateExpression: aws.String("SET #username = :username, #disabled = :disabled, #title = :title, #iconurl = :iconurl, #email = :email, #user_region = :user_region"),
		ReturnValues:     types.ReturnValueAllNew,
	})
	if err != nil {
		return nil, err
	}

	// If elo attribute is not set (new users) it is upserted with the BASEELO
	if _, ok := output.Attributes["elo"]; !ok {
		output, err = dynamoClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName: aws.String(tableName),
			Key: map[string]types.AttributeValue{
				"subject": &types.AttributeValueMemberS{Value: subject},
			},
			ExpressionAttributeNames: map[string]string{
				"#elo": "elo",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":elo": &types.AttributeValueMemberN{Value: baseElo},
			},
			UpdateExpression:    aws.String("SET #elo = :elo"),
			ConditionExpression: aws.String("attribute_not_exists(elo)"),
			ReturnValues:        types.ReturnValueAllNew,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to generate base elo: %v", err)
		}
	}

	var user UserOutput
	if err := attributevalue.UnmarshalMap(output.Attributes, &user); err != nil {
		return nil, fmt.Errorf("failed to deserialize upsert output")
	}

	return &user, nil
}
