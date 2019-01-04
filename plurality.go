package main

import (
//"fmt"
//"sync"
)

//PluralityMethod is a type of election method that can be used through the Method interface
type PluralityMethod struct {
	Electorate *Electorate       //reference to relevant electorate
	Winner     int               //index of winning candidate
	Ballots    []PluralityBallot //slice containing all ballots
	Utility    float64           //average utility per voter achieved by winning candidate
}

//Create creates the struct members needed to run the election
func (m *PluralityMethod) Create(e *Electorate) {
	m.Ballots = make([]PluralityBallot, len(e.Voters))
	m.Electorate = e
	m.Winner = -1
}

//GetWinner returns the index of the winning candidate
func (m *PluralityMethod) GetWinner() int {
	return m.Winner
}

//GetUtility returns the average utility for the winning candidate
func (m *PluralityMethod) GetUtility() float64 {
	return m.Utility
}

//Run creates ballots and tabulates the winner
func (m *PluralityMethod) Run() {

	for i := range m.Electorate.Voters {
		if m.Electorate.Voters[i].Strategic {
			m.Ballots[i] = m.VoteStrategic(&m.Electorate.Voters[i])
		} else {
			m.Ballots[i] = m.Vote(&m.Electorate.Voters[i])
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

	m.calcUtility()

	m.Ballots = nil
}

//calculates the average utility for the winning candidate
func (m *PluralityMethod) calcUtility() {
	u := 0.0
	for _, v := range m.Electorate.Voters {
		u += v.Utilities[m.Winner]
	}

	m.Utility = u / float64(len(m.Electorate.Voters))
}

//Vote creates a ballot for an honest voter
func (m *PluralityMethod) Vote(v *Voter) PluralityBallot {
	//an honest plurality voter just picks their favorite
	return PluralityBallot{Choice: findFavorite(v.Utilities)}
}

//VoteStrategic creates a ballot for a strategic voter
func (m *PluralityMethod) VoteStrategic(v *Voter) PluralityBallot {
	//a strategic plurality voter votes for their preferred major candidate
	iMax := 0
	uMax := 0.0
	for i, u := range v.Utilities {
		if u > uMax && m.Electorate.Candidates[i].Major {
			uMax = u
			iMax = i
		}
	}

	return PluralityBallot{Choice: iMax}
}

//PluralityBallot has a single field that holds the index of the chosen candidate
type PluralityBallot struct {
	Choice int //index of chosen candidate
}
