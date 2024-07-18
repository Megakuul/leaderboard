package query

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func FetchByUser(dynamoClient *dynamodb.Client, ctx context.Context, tableName string, username string) (*User, error) {
	output, err := dynamoClient.Query(ctx, &dynamodb.QueryInput{
		TableName: aws.String(tableName),
		IndexName: aws.String("username_gsi"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":username": &types.AttributeValueMemberS{Value: username},
		},
		KeyConditionExpression: aws.String("username = :username"),
		Limit:                  aws.Int32(1),
	})
	if err != nil {
		return nil, err
	}
	var users []User
	err = attributevalue.UnmarshalListOfMaps(output.Items, &users)
	if err != nil {
		return nil, err
	}
	if len(users) < 1 {
		return nil, fmt.Errorf("user not found")
	}
	return &users[0], nil
}
