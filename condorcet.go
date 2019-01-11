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

//finds the condorcet loser in a group of candidates, if one exists
func (e *Electorate) findGroupCondorcetLoser(candidates []int) (bool, int) {

Loop:
	for i := 0; i < len(candidates); i++ {
		for j := 0; j < len(candidates); j++ {
			if i == j {
				continue
			}

			if e.headToHead(candidates[i], candidates[j]) == true {
				continue Loop
			}
		}

		//a condorcet loser was found, so return true and the index of the loser
		return true, candidates[i]
	}

	return false, -1
}
