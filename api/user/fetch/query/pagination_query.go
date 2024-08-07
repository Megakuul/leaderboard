package query

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func FetchByPage(dynamoClient *dynamodb.Client, ctx context.Context, tableName string, pageSize int32, lastPageKey, region string) ([]UserOutput, string, error) {
	if pageSize > MAX_PAGESIZE {
		pageSize = MAX_PAGESIZE
	}

	var pageKey map[string]types.AttributeValue = nil
	if lastPageKey != "" {
		var err error
		pageKey, err = deserializePageKey(lastPageKey)
		if err != nil {
			return nil, "", fmt.Errorf("failed to deserialize lastPageKey: %v", err)
		}
	}

	output, err := dynamoClient.Query(ctx, &dynamodb.QueryInput{
		TableName: aws.String(tableName),
		IndexName: aws.String("region_gsi"),
		ExpressionAttributeNames: map[string]string{
			"#user_region": "user_region",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":user_region": &types.AttributeValueMemberS{Value: region},
		},
		KeyConditionExpression: aws.String("#user_region = :user_region"),
		Limit:                  aws.Int32(pageSize),
		ScanIndexForward:       aws.Bool(false),
		ExclusiveStartKey:      pageKey,
	})

	if err != nil {
		return nil, "", err
	}
	var users []UserOutput
	err = attributevalue.UnmarshalListOfMaps(output.Items, &users)
	if err != nil {
		return nil, "", err
	}
	if len(output.LastEvaluatedKey) < 1 {
		return users, "", nil
	}
	newPageKey, err := serializePageKey(output.LastEvaluatedKey)
	if err != nil {
		return nil, "", fmt.Errorf("failed to serialize new pagekey: %v", err)
	}
	return users, newPageKey, nil
}

func serializePageKey(pageKey map[string]types.AttributeValue) (string, error) {
	var translatedMap map[string]interface{}
	if err := attributevalue.UnmarshalMap(pageKey, &translatedMap); err != nil {
		return "", err
	}
	encodedMap, err := json.Marshal(&translatedMap)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(encodedMap), nil
}

func deserializePageKey(pageKey string) (map[string]types.AttributeValue, error) {
	decodedPageKey, err := base64.RawURLEncoding.DecodeString(pageKey)
	if err != nil {
		return nil, err
	}
	var decodedMap map[string]interface{}
	err = json.Unmarshal(decodedPageKey, &decodedMap)
	if err != nil {
		return nil, err
	}
	return attributevalue.MarshalMap(decodedMap)
}
