// contains wrappers for database query functions.
// main purpose is to abstract some boilerplate code
// away from the handler.
package query

type ParticipantOutput struct {
	Subject       string `dynamodbav:"subject"`
	Username      string `dynamodbav:"username"`
	EloUpdate     int    `dynamodbav:"elo_update"`
	Confirmed     bool   `dynamodbav:"confirmed"`
	ConfirmSecret string `dynamodbav:"confirm_secret"`
}

type GameOutput struct {
	GameId       string                       `dynamodbav:"gameid"`
	Readonly     bool                         `dynamodbav:"readonly"`
	Participants map[string]ParticipantOutput `dynamodbav:"participants"`
}
