package query

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func FetchByElo(dynamoClient *dynamodb.Client, ctx context.Context, tableName string, pageSize int32, region, elo string) ([]UserOutput, error) {
	if pageSize > MAX_PAGESIZE {
		pageSize = MAX_PAGESIZE
	}

	output, err := dynamoClient.Query(ctx, &dynamodb.QueryInput{
		TableName: aws.String(tableName),
		IndexName: aws.String("region_gsi"),
		ExpressionAttributeNames: map[string]string{
			"#user_region": "user_region",
			"#elo":         "elo",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":user_region": &types.AttributeValueMemberN{Value: region},
			":elo":         &types.AttributeValueMemberN{Value: elo},
		},
		KeyConditionExpression: aws.String("#user_region = :user_region AND #elo <= :elo"),
		ScanIndexForward:       aws.Bool(true),
		Limit:                  aws.Int32(pageSize),
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
