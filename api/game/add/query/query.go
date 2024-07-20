// contains wrappers for database query functions.
// main purpose is to abstract some boilerplate code
// away from the handler.
package query

type UserOutput struct {
	Subject  string `dynamodbav:"subject"`
	Username string `dynamodbav:"username"`
	Elo      int    `dynamodbav:"elo"`
	Email    string `dynamodbav:"email"`
}
