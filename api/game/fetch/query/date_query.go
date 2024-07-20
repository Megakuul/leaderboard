package query

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func FetchByDate(dynamoClient *dynamodb.Client, ctx context.Context, tableName string, date string) ([]GameOutput, error) {
	output, err := dynamoClient.Query(ctx, &dynamodb.QueryInput{
		TableName: aws.String(tableName),
		IndexName: aws.String("date_gsi"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":date": &types.AttributeValueMemberS{Value: date},
		},
		KeyConditionExpression: aws.String("date = :date"),
	})
	if err != nil {
		return nil, err
	}
	var games []GameOutput
	err = attributevalue.UnmarshalListOfMaps(output.Items, &games)
	if err != nil {
		return nil, err
	}
	return games, nil
}
