package main

import (
	//"fmt"
	"sync"
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
func (m *IRVMethod) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	//fmt.Println("creating IRV ballots")
	for i, v := range m.Electorate.Voters {
		if v.Strategic {
			m.Ballots[i] = m.VoteStrategic(v)
		} else {
			m.Ballots[i] = m.Vote(v)
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
}

//if there is a winner, returns (true, winner index), otherwise returns (false, last place index)
func (m *IRVMethod) checkForWinner() (bool, int) {
	//fmt.Println("checking for winner")
	leader := -1
	highVotes := 0
	loser := -1
	lowVotes := len(m.Ballots)

	for k, b := range m.Buckets {
		//fmt.Println(k, len(b))

		if len(b) > highVotes {
			leader = k
			highVotes = len(b)
		}

		if len(b) <= lowVotes {
			loser = k
			lowVotes = len(b)
		}
	}

	//if the leading candidate has a majority or there are only 2 candidates left, there is a winner
	if highVotes > len(m.Ballots)/2 || len(m.Buckets) == 2 {
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
	for _, b := range ballots {
		for {
			//increment choice on ballot
			b.LastChoice++

			//see if that candidate remains
			_, ok := m.Buckets[b.LastChoice]

			//if that candidate's index still exists, add ballot to bucket
			if ok {
				m.Buckets[b.Choices[b.LastChoice]] = append(m.Buckets[b.Choices[b.LastChoice]], b)
				break

				//if not, check to see if this ballot is expired and if so, discard ballot
			} else if b.LastChoice > len(m.Electorate.Candidates) {
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
func (m *IRVMethod) Vote(v Voter) IRVBallot {
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

	//fmt.Println(ballot)
	return ballot

}

//VoteStrategic creates a ballot for a strategic voter
func (m *IRVMethod) VoteStrategic(v Voter) IRVBallot {
	//strategic IRV voters will rank their preferred major candidate first and the other major candidate last

	return m.Vote(v)
}

//IRVBallot has a slice with indices of candidates in preferential order
type IRVBallot struct {
	Choices    []int //slice of candidate indices
	LastChoice int   //index of last choice read from this ballot
}
