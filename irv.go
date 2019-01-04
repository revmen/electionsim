package main

import (
//"fmt"
//"sync"
)

//IRVMethod is a type of election method that can be used through the Method interface
type IRVMethod struct {
	Electorate *Electorate         //reference to relevant electorate
	Winner     int                 //index of winning candidate
	Ballots    []IRVBallot         //slice containing all ballots
	Buckets    map[int][]IRVBallot //map of ballot slices used to tabulate results
	Utility    float64             //average utility per voter achieved by winning candidate
}

//Create creates the struct members needed to run the election
func (m *IRVMethod) Create(e *Electorate) {
	m.Ballots = make([]IRVBallot, len(e.Voters))
	m.Electorate = e
	m.Winner = -1

	//initialize each "bucket", which is a slice of ballots
	m.Buckets = make(map[int][]IRVBallot)
	for i := range m.Electorate.Candidates {
		m.Buckets[i] = make([]IRVBallot, 0)
	}
}

//GetWinner returns the index of the winning candidate
func (m *IRVMethod) GetWinner() int {
	return m.Winner
}

//GetUtility returns the average utility for the winning candidate
func (m *IRVMethod) GetUtility() float64 {
	return m.Utility
}

//Run creates ballots and tabulates the winner
func (m *IRVMethod) Run() {

	//fmt.Println("creating IRV ballots")
	for i := range m.Electorate.Voters {
		if m.Electorate.Voters[i].Strategic {
			m.Ballots[i] = m.VoteStrategic(&m.Electorate.Voters[i])
		} else {
			m.Ballots[i] = m.Vote(&m.Electorate.Voters[i])
		}
	}

	m.sortBallots(m.Ballots)

	//check for a winner
	//if no winner, eliminate last place and repeat
	var isWinner bool
	var ci int

	for {
		isWinner, ci = m.checkForWinner()

		if isWinner {
			break
		}

		m.eliminateCandidate(ci)
	}

	m.Winner = ci

	m.calcUtility()

	//m.Ballots = nil
}

//if there is a winner, returns (true, winner index), otherwise returns (false, last place index)
func (m *IRVMethod) checkForWinner() (bool, int) {
	//fmt.Println("checking for winner")
	leader := -1
	highVotes := 0
	loser := -1
	lowVotes := len(m.Ballots)

	for k := range m.Buckets {
		//fmt.Println(k, len(b))

		if len(m.Buckets[k]) > highVotes {
			leader = k
			highVotes = len(m.Buckets[k])

			//if the leading candidate has a majority, they are the winner
			if highVotes > len(m.Ballots)/2 {
				return true, leader
			}
		}

		if len(m.Buckets[k]) <= lowVotes {
			loser = k
			lowVotes = len(m.Buckets[k])
		}
	}

	//if there are only 2 candidates left, there is a winner
	if len(m.Buckets) == 2 {
		//fmt.Println("winner found", len(m.Buckets), len(m.Electorate.Candidates))
		return true, leader
	}

	//fmt.Println("no winner")
	return false, loser
}

//remove the indicated candidate and resort that candidate's ballots according to their next choice
func (m *IRVMethod) eliminateCandidate(i int) {
	//fmt.Println("eliminating candidate", i)
	bucket, ok := m.Buckets[i]
	if ok {
		//fmt.Println("key found")
		delete(m.Buckets, i)
		m.sortBallots(bucket)
	}

	//fmt.Println(m.Buckets)
}

//move each ballot in slice provided to the bucket of the next remaining candidate
func (m *IRVMethod) sortBallots(ballots []IRVBallot) {
	//fmt.Println("sorting ballots")
	for k := range ballots {
		for {
			//increment choice on ballot
			ballots[k].LastChoice++

			//see if that candidate remains
			_, ok := m.Buckets[ballots[k].LastChoice]

			//if that candidate's index still exists, add ballot to bucket
			if ok {
				m.Buckets[ballots[k].Choices[ballots[k].LastChoice]] = append(m.Buckets[ballots[k].Choices[ballots[k].LastChoice]], ballots[k])
				break

				//if not, check to see if this ballot is expired and if so, discard ballot
			} else if ballots[k].LastChoice > len(m.Electorate.Candidates) {
				break

				//otherwise, try next choice
			} else {
				continue
			}

		}
	}
	//fmt.Println("done sorting")
}

//calculates the average utility for the winning candidate
func (m *IRVMethod) calcUtility() {
	u := 0.0
	for _, v := range m.Electorate.Voters {
		u += v.Utilities[m.Winner]
	}

	m.Utility = u / float64(len(m.Electorate.Voters))
}

//Vote creates a ballot for an honest voter
func (m *IRVMethod) Vote(v *Voter) IRVBallot {
	ballot := IRVBallot{Choices: make([]int, 0), LastChoice: -1}

UtilitiesLoop:

	for i, u := range v.Utilities {

		for j, c := range ballot.Choices {

			if u > v.Utilities[c] {
				//insert i at j
				ballot.Choices = append(ballot.Choices, 0)
				copy(ballot.Choices[j+1:], ballot.Choices[j:])
				ballot.Choices[j] = i

				//next utility value
				continue UtilitiesLoop
			}
		}

		//if a lower value isn't found, append to bottom of Choices
		ballot.Choices = append(ballot.Choices, i)

	}

	//fmt.Println("honest", ballot)
	return ballot

}

//VoteStrategic creates a ballot for a strategic voter
func (m *IRVMethod) VoteStrategic(v *Voter) IRVBallot {
	//strategic IRV voters will rank their preferred major candidate first and the other major candidate last
	preferredMajor := findFavoriteMajor(v.Utilities, m.Electorate.Candidates)
	otherMajor := findOtherMajor(preferredMajor, m.Electorate.Candidates)

	ballot := IRVBallot{Choices: make([]int, 0), LastChoice: -1}

	//first choice is preferred major
	ballot.Choices = append(ballot.Choices, preferredMajor)

UtilitiesLoop:

	for i, u := range v.Utilities {

		//skip the Major candidates because we know where they go already
		if i == preferredMajor || i == otherMajor {
			continue
		}

		for j, c := range ballot.Choices {

			//skip the first spot because it's already set
			if j == 0 {
				continue
			}

			if u > v.Utilities[c] {
				//insert i at j
				ballot.Choices = append(ballot.Choices, 0)
				copy(ballot.Choices[j+1:], ballot.Choices[j:])
				ballot.Choices[j] = i

				//next utility value
				continue UtilitiesLoop
			}
		}

		//if a lower value isn't found, append to bottom of Choices
		ballot.Choices = append(ballot.Choices, i)

	}

	//put other major in last spot
	ballot.Choices = append(ballot.Choices, otherMajor)

	//fmt.Println("strategic", ballot)
	return ballot
}

//IRVBallot has a slice with indices of candidates in preferential order
type IRVBallot struct {
	Choices    []int //slice of candidate indices
	LastChoice int   //index of last choice read from this ballot
}
