package main

import (
//"sync"
)

//ApprovalMethod is a type of election method that can be used through the Method interface
type ApprovalMethod struct {
	Electorate *Electorate      //reference to relevant electorate
	Winner     int              //index of winning candidate
	Ballots    []ApprovalBallot //slice containing all ballots
	Utility    float64          //average utility per voter achieved by winning candidate
}

//Create creates the struct members needed to run the election
func (m *ApprovalMethod) Create(e *Electorate) {
	m.Ballots = make([]ApprovalBallot, len(e.Voters))
	m.Electorate = e
	m.Winner = -1
}

//GetWinner returns the index of the winning candidate
func (m *ApprovalMethod) GetWinner() int {
	return m.Winner
}

//GetUtility returns the average utility for the winning candidate
func (m *ApprovalMethod) GetUtility() float64 {
	return m.Utility
}

//Run creates ballots and tabulates the winner
func (m *ApprovalMethod) Run() {

	for i := range m.Electorate.Voters {
		if m.Electorate.Voters[i].Strategic {
			m.Ballots[i] = m.VoteStrategic(&m.Electorate.Voters[i])
		} else {
			m.Ballots[i] = m.Vote(&m.Electorate.Voters[i])
		}
	}

	votes := make([]int, len(m.Electorate.Candidates))

	for i := range m.Ballots {
		for j := range m.Ballots[i].Approvals {
			if m.Ballots[i].Approvals[j] {
				votes[j]++
			}
		}
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
func (m *ApprovalMethod) calcUtility() {
	u := 0.0
	for _, v := range m.Electorate.Voters {
		u += v.Utilities[m.Winner]
	}

	m.Utility = u / float64(len(m.Electorate.Voters))
}

//Vote creates a ballot for an honest voter
func (m *ApprovalMethod) Vote(v *Voter) ApprovalBallot {
	ballot := ApprovalBallot{Approvals: make([]bool, len(v.Utilities))}

	for i, u := range v.Utilities {
		if u > v.ApprovalThreshold {
			ballot.Approvals[i] = true
		} else {
			ballot.Approvals[i] = false
		}
	}

	return ballot

}

//VoteStrategic creates a ballot for a strategic voter
func (m *ApprovalMethod) VoteStrategic(v *Voter) ApprovalBallot {
	//strategic approval voters will bullet vote if their preferred candidate is a major
	//if their favorite is not a major, they will vote honestly

	favorite := findFavorite(v.Utilities)
	ballot := ApprovalBallot{Approvals: make([]bool, len(v.Utilities))}

	if !m.Electorate.Candidates[favorite].Major {
		ballot = m.Vote(v)
	} else {
		for i := range v.Utilities {
			if i == favorite {
				ballot.Approvals[i] = true
			} else {
				ballot.Approvals[i] = false
			}
		}

	}
	return ballot
}

//ApprovalBallot has boolean for each candidate
type ApprovalBallot struct {
	Approvals []bool //slice of bools indicating up or down votes for each candidate
}
