package main

import (
//"fmt"
)

//ApprovalMethod is a type of election method that can be used through the Method interface
type ApprovalMethod struct {
	Electorate *Electorate
	Winner     int //index of winning candidate
	Ballots    []ApprovalBallot
}

//Create creates the struct members needed to run the election
func (m *ApprovalMethod) Create(e *Electorate) {
	m.Ballots = make([]ApprovalBallot, 0)
	m.Electorate = e
	m.Winner = -1
}

//GetWinner returns the index of the winning candidate
func (m *ApprovalMethod) GetWinner() int {
	return m.Winner
}

//Run creates ballots and tabulates the winner
func (m *ApprovalMethod) Run() {

	for _, v := range m.Electorate.Voters {
		if v.Strategic {
			m.VoteStrategic(v)
		} else {
			m.Vote(v)
		}
	}

	votes := make([]int, len(m.Electorate.Candidates))

	for i := range m.Ballots {
		for j, a := range m.Ballots[i].Approvals {
			if a {
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
}

//Vote creates a ballot for an honest voter
func (m *ApprovalMethod) Vote(v Voter) {
	ballot := ApprovalBallot{Approvals: make([]bool, len(v.Utilities))}

	for i, u := range v.Utilities {
		if u > v.ApprovalThreshold {
			ballot.Approvals[i] = true
		} else {
			ballot.Approvals[i] = false
		}
	}

	m.Ballots = append(m.Ballots, ballot)

}

//VoteStrategic creates a ballot for a strategic voter
func (m *ApprovalMethod) VoteStrategic(v Voter) {
	//strategic approval voters will bullet vote if their preferred candidate is a major
	//if their favorite is not a major, they will vote honestly

	favorite := findFavorite(v.Utilities)

	if !m.Electorate.Candidates[favorite].Major {
		m.Vote(v)
	} else {

		ballot := ApprovalBallot{Approvals: make([]bool, len(v.Utilities))}

		for i := range v.Utilities {
			if i == favorite {
				ballot.Approvals[i] = true
			} else {
				ballot.Approvals[i] = false
			}
		}

		m.Ballots = append(m.Ballots, ballot)

	}
}

//ApprovalBallot has boolean for each candidate
type ApprovalBallot struct {
	Approvals []bool //slice of bools indicating up or down votes for each candidate
}
