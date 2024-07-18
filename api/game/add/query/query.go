// contains wrappers for fetch query functions.
// main purpose is to abstract some boilerplate code
// away from the handler.
package query

type User struct {
	Username string `dynamodbav:"username"`
	Elo      string `dynamodbav:"elo"`
	Email    string `dynamodbav:"email"`
}
