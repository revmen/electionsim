package main

import "math"

// ScoreMethod : Each voter gives each candidate a score within some range specified by `min` and `max`.
// The winner is the candidate with the most votes.
type ScoreMethod struct {
	min int
	max int
}

// NewScoreMethod is the "constructor" for ScoreMethod. It requires a minimum and maximum score to be specified.
func NewScoreMethod(min, max int) ScoreMethod {
	return ScoreMethod{min, max}
}

// NewAdaptedScoreMethod is a convinience function to construct a ScoreMethod and adapt it to the normal Method interface.
func NewAdaptedScoreMethod(min, max int) AdaptedMethod {
	ScoreMethod := NewScoreMethod(min, max)
	return AdaptSimpleMethod(&ScoreMethod)
}

// FindWinner finds the index of the Score winner of the provided Electorate.
func (m *ScoreMethod) FindWinner(electorate *Electorate) int {
	sums := make([]int, len(electorate.Candidates))

	for _, voter := range electorate.Voters {
		favoriteMajor := findFavoriteMajor(voter.Utilities, electorate.Candidates)
		theshold := voter.Utilities[favoriteMajor]
		for j, score := range m.vote(&voter, theshold) {
			sums[j] += score
		}
	}

	return findLargestIndex(sums)
}

func findLargestIndex(list []int) int {
	largestIndex := 0
	largest := list[largestIndex]

	for i, value := range list {
		if value > largest {
			largest = value
			largestIndex = i
		}
	}

	return largestIndex
}

func (m *ScoreMethod) vote(voter *Voter, strategicThreshold float64) []int {

	if voter.Strategic {
		return thresholdClamp(voter.Utilities, strategicThreshold, m.min, m.max)
	}

	return linearScale(voter.Utilities, m.min, m.max)
}

// linearScale returns a copy of "list" with its values linearly scaled such that its smallest value becomes "min" and its largest value becomes "max"
// If all values in "list" are the same, the final values will all be zero instead
func linearScale(list []float64, min, max int) []int {
	result := make([]int, len(list))

	smallest := list[0]
	largest := list[0]

	for _, util := range list {
		if util > largest {
			largest = util
		} else if util < smallest {
			smallest = util
		}
	}

	if largest == smallest {
		return result
	}

	scalefactor := float64(max-min) / (largest - smallest)

	for i, util := range list {
		// smallest shifted to zero, which isn't effected by scaling (stays 0), then min is added (ends up at min)
		// largest shifted down, then scaled to (max-min), then min is added (ends up at max)
		shifted := util - smallest
		scaled := shifted * scalefactor
		rounded := int(math.Floor(scaled + 0.5))
		result[i] = rounded + min
	}

	return result
}

// thresholdClamp returns a copy of "list" with its values clamped such that values below a specified "threshold" become "min" and the remainder become "max".
// Values equal to "threshold" become "max".
func thresholdClamp(list []float64, threshold float64, min, max int) []int {

	clamped := make([]int, len(list))

	for i, util := range list {
		if util >= threshold {
			clamped[i] = max
		} else {
			clamped[i] = min
		}
	}

	return clamped
}
