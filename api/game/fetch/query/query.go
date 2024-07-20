// contains wrappers for fetch query functions.
// each query is abstracted in its own function as they utilize different
// dynamodb tools (indexes, etc.)
package query

type ParticipantOutput struct {
	Username  string `dynamodbav:"username" json:"username"`
	Placement int    `dynamodbav:"placement" json:"placement"`
	Points    int    `dynamodbav:"points" json:"points"`
	Elo       int    `dynamodbav:"elo" json:"elo"`
	EloUpdate int    `dynamodbav:"elo_update" json:"elo_update"`
	Confirmed bool   `dynamodbav:"confirmed" json:"confirmed"`
}

type GameOutput struct {
	GameId       string              `dynamodbav:"gameid" json:"gameid"`
	Date         string              `dynamodbav:"date" json:"date"`
	Participants []ParticipantOutput `dynamodbav:"participants" json:"participants"`
}
