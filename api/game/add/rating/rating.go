// contains functions to calculate rating updates
package rating

import (
	"sort"

	"github.com/megakuul/leaderboard/api/game/add/query"
)

type ParticipantInput struct {
	UserRef   *query.UserOutput
	Team      int
	Rating    int
	Points    int
	Placement int
}

type team struct {
	Participants []*ParticipantInput
	Rating       int
	Points       int
}

type ParticipantOutput struct {
	UserRef      *query.UserOutput
	RatingUpdate int
	Team         int
	Rating       int
	Points       int
	Placement    int
}

// CalculateRatingUpdate uses a simple algorithm to update the rating for all participants of a game.
// It operates in four fundamental steps:
// 1. Reverse sort participants by placement and increment points by placementPoints * index.
// 2. Add participants to teams (based on the participants "team" tag) and sum all the participants rating and points per team.
// 3. Create a hypothesis & evidence value for every team. This means simply calculating the percentage of rating & points based on the current game.
// 4. Calculate the difference between hypothesis & evidence multiply it with maxLossNumber and setting it as rating update on all participants of the team.
func CalculateRatingUpdate(participants []ParticipantInput, placementPoints int, maxLossNumber int) []ParticipantOutput {
	// Reverse sort, to assign points based on index position
	sort.Slice(participants, func(i, j int) bool {
		return participants[i].Placement < participants[j].Placement
	})

	// teams represent a intermediate calculation entity.
	// They are used to ensure all players of one team have the same rating update.
	teams := map[int]team{}

	// Rating is used for hypothesis calculation
	var combinedRating int
	// Points are used for evidence calculation
	var combinedPoints int

	// In one iteration 3 things are done:
	// add placement points, calculate combined rating + points and add participant to calculationEntity
	for i, part := range participants {
		// Step 1. add placement points to the participants points
		part.Points += i * placementPoints

		// Step 2. add participant to combinedRating and combinedPoints
		combinedPoints += part.Points
		combinedRating += part.Rating

		// Step 3. add the participant to a calculationEntity
		entity, ok := teams[part.Team]
		if ok {
			entity.Participants = append(entity.Participants, &part)
			entity.Rating += part.Rating
			entity.Points += part.Points
		} else {
			teams[part.Team] = team{
				Participants: []*ParticipantInput{&part},
				Rating:       part.Rating,
				Points:       part.Points,
			}
		}
	}

	outputParticipants := []ParticipantOutput{}
	for _, entity := range teams {
		// hypothesis is the percentage of rating in this game
		hypothesis := entity.Rating / combinedRating
		// evidence is the percentage of points in this game
		evidence := entity.Points / combinedPoints

		// update is the difference between hypothesis and evidence converted to elo with the maxLossNumber
		update := maxLossNumber * (evidence - hypothesis)

		for _, part := range entity.Participants {
			outputParticipants = append(outputParticipants, ParticipantOutput{
				UserRef:      part.UserRef,
				RatingUpdate: update,
				Team:         part.Team,
				Rating:       part.Rating,
				Points:       part.Points,
				Placement:    part.Placement,
			})
		}
	}
	return outputParticipants
}
