package main

import (
	"fmt"
)

//PluralityMethod is a type of election method that can be used through the Method interface
type PluralityMethod struct {
	Electorate *Electorate
	Winner     int //index of winning candidate
	Ballots    []PluralityBallot
}

//Create creates the struct members needed to run the election
func (m *PluralityMethod) Create(e *Electorate) {
	m.Ballots = make([]PluralityBallot, 0)
	m.Electorate = e
	m.Winner = -1
}

//GetWinner returns the index of the winning candidate
func (m *PluralityMethod) GetWinner() int {
	return m.Winner
}

//Run creates ballots and tabulates the winner
func (m *PluralityMethod) Run() {

	for _, v := range m.Electorate.Voters {
		if v.Strategic {
			m.VoteStrategic(v)
		} else {
			m.Vote(v)
		}
	}

	votes := make([]int, len(m.Electorate.Candidates))

	for _, b := range m.Ballots {
		votes[b.Choice]++
	}

	winningVotes := 0

	for i := range votes {
		if votes[i] > winningVotes {
			winningVotes = votes[i]
			m.Winner = i
		}
	}
}

//Vote creates a ballot for an honest voter
func (m *PluralityMethod) Vote(v Voter) {
	//an honest plurality voter just picks their favorite
	m.Ballots = append(m.Ballots, PluralityBallot{Choice: findFavorite(v.Utilities)})
}

//VoteStrategic creates a ballot for a strategic voter
func (m *PluralityMethod) VoteStrategic(v Voter) {
	//a strategic plurality voter votes for their preferred major candidate
	iMax := 0
	uMax := 0.0
	for i, u := range v.Utilities {
		if u > uMax && m.Electorate.Candidates[i].Major {
			uMax = u
			iMax = i
		}
	}

	m.Ballots = append(m.Ballots, PluralityBallot{Choice: iMax})
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
