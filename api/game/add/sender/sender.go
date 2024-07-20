// contains wrapper for sending confirmation mail to aws ses.
package sender

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

type EmailConfirmRequest struct {
	Email     string
	Secret    string
	Placement int
	Points    int
	EloUpdate int
}

type emailTemplateInput struct {
	GameId    string `json:"gameid"`
	Secret    string `json:"secret"`
	Placement int    `json:"placement"`
	Points    int    `json:"points"`
	EloUpdate int    `json:"elo_update"`
}

func SendConfirmMails(sesClient *sesv2.Client, ctx context.Context, senderMail, mailTemplate, gameId string, emailRequests []EmailConfirmRequest) error {
	emailDestinations := []types.BulkEmailEntry{}
	for _, request := range emailRequests {
		templateInput := emailTemplateInput{
			GameId:    gameId,
			Secret:    request.Secret,
			Placement: request.Placement,
			Points:    request.Points,
			EloUpdate: request.EloUpdate,
		}
		templateInputSerialized, err := json.Marshal(&templateInput)
		if err != nil {
			return fmt.Errorf("failed to serialize mail input")
		}
		emailDestinations = append(emailDestinations, types.BulkEmailEntry{
			Destination: &types.Destination{
				ToAddresses: []string{
					request.Email,
				},
			},
			ReplacementEmailContent: &types.ReplacementEmailContent{
				ReplacementTemplate: &types.ReplacementTemplate{
					ReplacementTemplateData: aws.String(string(templateInputSerialized)),
				},
			},
		})
	}

	_, err := sesClient.SendBulkEmail(ctx, &sesv2.SendBulkEmailInput{
		BulkEmailEntries: emailDestinations,
		FromEmailAddress: aws.String(senderMail),
		DefaultContent: &types.BulkEmailContent{
			Template: &types.Template{
				TemplateData: aws.String("{}"),
				TemplateName: aws.String(mailTemplate),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to send templated email to SES")
	}
	return nil
}
