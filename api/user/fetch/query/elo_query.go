package query

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func FetchByElo(dynamoClient *dynamodb.Client, ctx context.Context, tableName string, elo string) ([]UserOutput, error) {
	output, err := dynamoClient.Query(ctx, &dynamodb.QueryInput{
		TableName: aws.String(tableName),
		IndexName: aws.String("elo_gsi"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":elo": &types.AttributeValueMemberS{Value: elo},
		},
		KeyConditionExpression: aws.String("elo = :elo"),
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
