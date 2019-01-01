package main

import (
//"fmt"
)

func findCondorcetWinner(e Electorate) int {
	winner := -1

Loop:
	for i := range e.Candidates {
		for j := range e.Candidates {
			if i == j {
				continue
			}

			if headToHead(e, i, j) == false {
				continue Loop
			}
		}

		winner = i
		break
	}

	return winner
}

// true if the candidate at index 1 (i1) beats the candidate at index 2 (i2) in a head-to-head matchup
func headToHead(e Electorate, i1, i2 int) bool {
	votes := 0

	for _, v := range e.Voters {
		if v.Utilities[i1] > v.Utilities[i2] {
			votes++
		}
	}

	return votes > len(e.Voters)/2
}
