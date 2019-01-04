package main

import ()

//identifies the condorcet winner, if any, for an electorate
func (e *Electorate) findCondorcetWinner() {
	winner := -1

Loop:
	for i := range e.Candidates {
		for j := range e.Candidates {
			if i == j {
				continue
			}

			if e.headToHead(i, j) == false {
				continue Loop
			}
		}

		winner = i
		break
	}

	e.CondorcetWinner = winner
}

// true if the candidate at index 1 (i1) beats the candidate at index 2 (i2) in a head-to-head matchup
func (e *Electorate) headToHead(i1, i2 int) bool {
	votes := 0

	for _, v := range e.Voters {
		if v.Utilities[i1] > v.Utilities[i2] {
			votes++
		}
	}

	return votes > len(e.Voters)/2
}
