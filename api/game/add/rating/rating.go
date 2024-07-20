// contains functions to calculate rating updates
package rating

import "github.com/megakuul/leaderboard/api/game/add/query"

type ParticipantInput struct {
	UserRef   *query.UserOutput
	Rating    int
	Points    int
	Placement int
}

type ParticipantOutput struct {
	UserRef      *query.UserOutput
	RatingUpdate int
	Rating       int
	Points       int
	Placement    int
}

func CalculateRatingUpdate(participants []ParticipantInput) []ParticipantOutput {

}
