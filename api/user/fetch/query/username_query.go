package query

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func FetchByUsername(dynamoClient *dynamodb.Client, ctx context.Context, tableName string, username string) ([]UserOutput, error) {
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
	var users []UserOutput
	err = attributevalue.UnmarshalListOfMaps(output.Items, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}
