// contains functions to calculate rating updates
package rating

import (
	"sort"

	"github.com/megakuul/leaderboard/api/game/add/query"
)

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

func CalculateRatingUpdate(participants []ParticipantInput, placementPoints int, maxLossNumber int) []ParticipantOutput {
	// Reverse sort, to assign points based on index position
	sort.Slice(participants, func(i, j int) bool {
		return participants[i].Placement < participants[j].Placement
	})

	// Rating is used for hypothesis calculation
	var combinedRating int
	// Points are used for evidence calculation
	var combinedPoints int

	// In one iteration 3 things are done: add placement points, calculate rating and points of all players
	for i, part := range participants {
		part.Points += i * placementPoints
		combinedPoints += part.Points
		combinedRating += part.Rating
	}

	outputParticipants := []ParticipantOutput{}
	for _, part := range participants {
		// hypothesis is the percentage of rating inside the participant pool
		hypothesis := part.Rating / combinedRating
		// evidence is the percentage of points inside the participant pool
		evidence := part.Points / combinedPoints

		update := maxLossNumber * (evidence - hypothesis)

		outputParticipants = append(outputParticipants, ParticipantOutput{
			UserRef:      part.UserRef,
			RatingUpdate: update,
			Rating:       part.Rating,
			Points:       part.Points,
			Placement:    part.Placement,
		})
	}
	return outputParticipants
}
