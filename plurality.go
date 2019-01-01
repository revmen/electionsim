package main

import (
	"fmt"
)

//PluralityElection is a type of election
type PluralityElection struct {
	Winner  int //index of winning candidate
	Ballots []PluralityBallot
}

//DoElection creates ballots and tabulates the winner
func (pe *PluralityElection) DoElection(e Electorate) {
	pe.Ballots = make([]PluralityBallot, 0)

	for _, v := range e.Voters {
		if v.Strategic {
			pe.VoteStrategic(v, e)
		} else {
			pe.Vote(v)
		}
	}

	votes := make([]int, len(e.Candidates))

	for _, b := range pe.Ballots {
		votes[b.Choice]++
	}

	winningVotes := 0

	for i := range votes {
		if votes[i] > winningVotes {
			winningVotes = votes[i]
			pe.Winner = i
		}
	}
}

//Vote creates a ballot for a voter that is not strategic
func (pe *PluralityElection) Vote(v Voter) {
	iMax := 0
	uMax := 0.0
	for i, u := range v.Utilities {
		if u > uMax {
			uMax = u
			iMax = i
		}
	}

	pe.Ballots = append(pe.Ballots, PluralityBallot{Choice: iMax})

}

//VoteStrategic creates a ballot for a voter that always votes for a major candidate
func (pe *PluralityElection) VoteStrategic(v Voter, e Electorate) {
	iMax := 0
	uMax := 0.0
	for i, u := range v.Utilities {
		if u > uMax && e.Candidates[i].Major {
			uMax = u
			iMax = i
		}
	}

	pe.Ballots = append(pe.Ballots, PluralityBallot{Choice: iMax})
}

//PluralityBallot has a single field that holds the index of the chosen candidate
type PluralityBallot struct {
	Choice int //index of chosen candidate
}

//for troubleshooting
func printPluralityVotes(e Electorate) {
	votes := make([]int, len(e.Candidates))

	for _, v := range e.Voters {
		iMax := 0
		uMax := 0.0
		for i, u := range v.Utilities {
			if u > uMax {
				uMax = u
				iMax = i
			}
		}

		votes[iMax]++
	}

	for i, c := range e.Candidates {
		//fmt.Printf("%# v", pretty.Formatter(e))
		fmt.Println(c.Name, votes[i])
	}
}
