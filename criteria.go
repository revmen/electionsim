package main

import (
//"fmt"
)

func findCondorcetWinner(e Electorate) string {
	winnerName := "no condorcet winner"

Loop:
	for i1, c1 := range e.Candidates {
		for i2, c2 := range e.Candidates {
			if c1.Name == c2.Name {
				continue
			}

			if headToHead(e, i1, i2) == false {
				continue Loop
			}
		}

		winnerName = c1.Name
	}

	return winnerName
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
