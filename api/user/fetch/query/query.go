// contains wrappers for database query functions.
// each query is abstracted in its own function as they utilize different
// dynamodb tools (indexes, pagination etc.)
package query

type UserOutput struct {
	Username string `dynamodbav:"username" json:"username"`
	Title    string `dynamodbav:"title" json:"title"`
	IconUrl  string `dynamodbav:"iconurl" json:"iconurl"`
	Elo      int    `dynamodbav:"elo" json:"elo"`
}
