// contains wrappers for fetch query functions.
// each query is abstracted in its own function as they utilize different
// dynamodb tools (indexes, pagination etc.)
package query

type User struct {
	Username string `dynamodbav:"username" json:"username"`
	Title    string `dynamodbav:"title" json:"title"`
	IconUrl  string `dynamodbav:"iconurl" json:"iconurl"`
	Elo      int    `dynamodbav:"elo" json:"elo"`
}
