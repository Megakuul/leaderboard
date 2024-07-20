package query

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func FetchById(dynamoClient *dynamodb.Client, ctx context.Context, tableName string, gameid string) (*GameOutput, error) {
	output, err := dynamoClient.Query(ctx, &dynamodb.QueryInput{
		TableName: aws.String(tableName),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":gameid": &types.AttributeValueMemberS{Value: gameid},
		},
		KeyConditionExpression: aws.String("gameid = :gameid"),
		Limit:                  aws.Int32(1),
	})
	if err != nil {
		return nil, err
	}
	var games []GameOutput
	err = attributevalue.UnmarshalListOfMaps(output.Items, &games)
	if err != nil {
		return nil, err
	}
	if len(games) < 1 {
		return nil, fmt.Errorf("game not found")
	}
	return &games[0], nil
}
