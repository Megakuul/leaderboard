// contains functions to calculate rating updates
package rating

import (
	"math"
	"sort"

	"github.com/megakuul/leaderboard/api/game/add/query"
)

const (
	UNDERDOG_BONUS_MULTIPLICATOR = 1
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
	Underdog     bool
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
func CalculateRatingUpdate(participants []ParticipantInput, placementPoints int, maxLossNumber int) []*ParticipantOutput {
	// Reverse sort, to assign points based on index position
	sort.Slice(participants, func(i, j int) bool {
		return participants[i].Placement < participants[j].Placement
	})

	// teams represent a intermediate calculation entity.
	// They are used to ensure all players of one team have the same rating update.
	teams := map[int]*team{}

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
			teams[part.Team] = &team{
				Participants: []*ParticipantInput{&part},
				Rating:       part.Rating,
				Points:       part.Points,
			}
		}
	}

	// this algorithm has a problem: when performing the hypothesis & evidence division
	// this yields a float64. As this elo system does only use integers, we need to convert the float64 back to int.
	// problem is that the elo system should not leak elo. for that reason, we need to put the remainders of divisons anywhere.
	// this is where the underdog commes into play, it is ref to the participant with the largest positive difference in the game.
	// at the end of the calculations the underdog rating bonus is added to the participants rating update.
	// all remainders are added to this underdog rating bonus in order to prevent elo leaking,
	// because this can lead to a negative underdog bonus there is a UNDERDOG_BONUS_MULTIPLICATOR constant, which is removed from every teams rating update
	// and added to the underdogRatingBonus. this increases the value of the underdog bonus and heavily reduces the chance of a negative underdog bonus.
	var underdogRef *ParticipantOutput = nil
	var underdogDifference float64 = 0.0
	underdogRatingBonus := float64(len(teams) * UNDERDOG_BONUS_MULTIPLICATOR)

	outputParticipants := []*ParticipantOutput{}
	for _, entity := range teams {
		setUnderdog := false

		// hypothesis is the percentage of rating in this game
		hypothesis := float64(entity.Rating) / float64(combinedRating)
		// evidence is the percentage of points in this game
		evidence := float64(entity.Points) / float64(combinedPoints)

		// calculate the difference, if the difference is positive and larger then the previous
		// underdogDiff, the underdog flag is set for this team.
		difference := evidence - hypothesis
		if difference > 0 && difference > underdogDifference {
			underdogDifference = difference
			setUnderdog = true
		}

		// underdog bonus multiplicator is removed
		baseUpdate := (float64(maxLossNumber) * (evidence - hypothesis)) - UNDERDOG_BONUS_MULTIPLICATOR
		// integer frac of the update is used for further calculations.
		updateNum := int(baseUpdate)
		// remaining float frac is shifted to the underdog bonus as we don't want to leak this.
		underdogRatingBonus += baseUpdate - float64(updateNum)

		// acquire the rating update per participant.
		// larger teams get smaller individual updates, as each member has less game impact.
		individualUpdate := updateNum / len(entity.Participants)

		// add the remainder of the update split per participant to the underdog bonus.
		underdogRatingBonus += float64(updateNum % len(entity.Participants))

		// flag to track the highest points reached in this team.
		maxPoints := 0
		for _, part := range entity.Participants {
			output := ParticipantOutput{
				UserRef:      part.UserRef,
				Underdog:     false,
				RatingUpdate: individualUpdate,
				Team:         part.Team,
				Rating:       part.Rating,
				Points:       part.Points,
				Placement:    part.Placement,
			}
			outputParticipants = append(outputParticipants, &output)

			// If underdog flag is set AND the participant has the most points of the team
			// he is set as underdogRef (dough this can change later on).
			if setUnderdog && part.Points > maxPoints {
				maxPoints = part.Points
				underdogRef = &output
			}
		}
	}

	// add underdog bonus. as we added all remainders to this, it should be an almost exact integer.
	if underdogRef != nil {
		underdogRef.RatingUpdate += int(math.Round(underdogRatingBonus))
		underdogRef.Underdog = true
	}

	return outputParticipants
}
