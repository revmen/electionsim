package main

import (
	"math"
)

func findUtilityWinner(e Electorate) (int, float64) {
	numVoters := len(e.Voters)

	winner := -1
	var winnerUtil float64

	var util float64

	for i := range e.Candidates {
		util = 0.0

		for _, v := range e.Voters {
			util += v.Utilities[i]
		}

		util = util / float64(numVoters)
		if util > winnerUtil {
			winner = i
			winnerUtil = util
		}

	}

	return winner, winnerUtil
}

//calculates the utilty for a voter from an elected candidate based on their distance in ideological space
func utility(v Voter, c Candidate) float64 {

	//max distance based on number of axes
	maxDistance := math.Sqrt(float64(len(v.Alignments)))

	//distance between voter and candidate
	d := distance(v.Alignments, c.Alignments)

	//normalize to 0.0 <-> 1.0
	d = d / maxDistance

	//reverse value so that bigger is better
	return 1 - d
}

//the direct distance between two sets of alignments
func distance(a1, a2 []float64) float64 {
	numAxes := len(a1)

	d := 0.0

	for i := 0; i < numAxes; i++ {
		d += math.Pow(a1[i]-a2[i], 2.0)
	}

	return math.Sqrt(d)
}
