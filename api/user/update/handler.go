package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type User struct {
	Username string `json:"username"`
	Title    string `json:"title"`
	Email    string `json:"email"`
	IconUrl  string `json:"iconurl"`
	Elo      string `json:"elo"`
}

type UpdateResponse struct {
	Message     string `json:"message"`
	UpdatedUser User   `json:"updated_user"`
}

func UpdateHandler(dynamoClient *dynamodb.Client) func(context.Context, events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return func(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		response, code, err := runUpdatehandler(dynamoClient, &request, ctx)
		if err != nil {
			return events.APIGatewayV2HTTPResponse{
				StatusCode: code,
				Headers:    map[string]string{"Content-Type": "text/plain"},
				Body:       err.Error(),
			}, nil
		}
		serializedResponse, err := json.Marshal(&response)
		if err != nil {
			return events.APIGatewayV2HTTPResponse{
				StatusCode: http.StatusInternalServerError,
				Headers:    map[string]string{"Content-Type": "text/plain"},
				Body:       "failed to serialize response",
			}, nil
		}
		return events.APIGatewayV2HTTPResponse{
			StatusCode: code,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       string(serializedResponse),
		}, nil
	}
}

func runUpdatehandler(dynamoClient *dynamodb.Client, request *events.APIGatewayV2HTTPRequest, ctx context.Context) (*UpdateResponse, int, error) {

	sub := request.RequestContext.Authorizer.JWT.Claims["sub"]
	if sub == "" {
		return nil, http.StatusUnprocessableEntity, fmt.Errorf("invalid sub claim in the ID token")
	}

	output, err := dynamoClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(USERTABLE),
		Key: map[string]types.AttributeValue{
			"subject": &types.AttributeValueMemberS{Value: sub},
		},
		ExpressionAttributeNames: map[string]string{
			"#username": "username",
			"#title":    "title",
			"#iconurl":  "iconurl",
			"#email":    "email",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":username": &types.AttributeValueMemberS{Value: request.RequestContext.Authorizer.JWT.Claims["preferred_username"]},
			":title":    &types.AttributeValueMemberS{Value: request.RequestContext.Authorizer.JWT.Claims["nickname"]},
			":iconurl":  &types.AttributeValueMemberS{Value: request.RequestContext.Authorizer.JWT.Claims["picture"]},
			":email":    &types.AttributeValueMemberS{Value: request.RequestContext.Authorizer.JWT.Claims["email"]},
		},
		UpdateExpression: aws.String("SET #username = :username, #title = :title, #iconurl = :iconurl, #email = :email"),
		ReturnValues:     types.ReturnValueAllNew,
	})
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("failed to update user: %v", err)
	}

	// If elo attribute is not set (new users) it is upserted with the BASEELO
	if _, ok := output.Attributes["elo"]; !ok {
		output, err = dynamoClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName: aws.String(USERTABLE),
			Key: map[string]types.AttributeValue{
				"subject": &types.AttributeValueMemberS{Value: sub},
			},
			ExpressionAttributeNames: map[string]string{
				"#elo": "elo",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":elo": &types.AttributeValueMemberN{Value: BASEELO},
			},
			UpdateExpression:    aws.String("SET #elo = :elo"),
			ConditionExpression: aws.String("attribute_not_exists(elo)"),
			ReturnValues:        types.ReturnValueAllNew,
		})
		if err != nil {
			return nil, http.StatusBadRequest, fmt.Errorf("failed to generate base elo for user: %v", err)
		}
	}

	username := ""
	if usernameAttr, ok := output.Attributes["username"].(*types.AttributeValueMemberS); ok {
		username = usernameAttr.Value
	}
	title := ""
	if titleAttr, ok := output.Attributes["title"].(*types.AttributeValueMemberS); ok {
		title = titleAttr.Value
	}
	email := ""
	if emailAttr, ok := output.Attributes["email"].(*types.AttributeValueMemberS); ok {
		email = emailAttr.Value
	}
	iconurl := ""
	if iconurlAttr, ok := output.Attributes["iconurl"].(*types.AttributeValueMemberS); ok {
		iconurl = iconurlAttr.Value
	}
	elo := ""
	if eloAttr, ok := output.Attributes["elo"].(*types.AttributeValueMemberN); ok {
		elo = eloAttr.Value
	}

	return &UpdateResponse{
		Message: "successfully updated user",
		UpdatedUser: User{
			Username: username,
			Title:    title,
			Email:    email,
			IconUrl:  iconurl,
			Elo:      elo,
		},
	}, http.StatusOK, nil
}
