// contains wrappers for database query functions.
// each query is abstracted in its own function as they utilize different
// dynamodb tools (indexes, pagination etc.)
package query

const (
	MAX_PAGESIZE = 100
)

type UserOutput struct {
	Username string `dynamodbav:"username" json:"username"`
	Disabled bool   `dynamodbav:"disabled" json:"disabled"`
	Region   string `dynamodbav:"user_region" json:"region"`
	Title    string `dynamodbav:"title" json:"title"`
	IconUrl  string `dynamodbav:"iconurl" json:"iconurl"`
	Elo      int    `dynamodbav:"elo" json:"elo"`
}
