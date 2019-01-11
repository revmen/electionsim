package main

import ()

//PairwiseElimMethod is a type of election method that can be used through the Method interface
type PairwiseElimMethod struct {
	Electorate *Electorate                  //reference to relevant electorate
	Winner     int                          //index of winning candidate
	Ballots    []PairwiseElimBallot         //slice containing all ballots
	Utility    float64                      //average utility per voter achieved by winning candidate
	Buckets    map[int][]PairwiseElimBallot //buckets of hate ballots for each candidate
}

//PairwiseElimBallot has a slice with indices of candidates in preferential order. It's identical to an IRV ballot
type PairwiseElimBallot struct {
	Choices    []int //slice of candidate indices
	LastChoice int   //index of last choice read from this ballot
}

//Create creates the struct members needed to run the election
func (m *PairwiseElimMethod) Create(e *Electorate) {
	m.Ballots = make([]PairwiseElimBallot, len(e.Voters))
	m.Electorate = e
	m.Winner = -1

	//initialize each "bucket", which is a slice of ballots
	m.Buckets = make(map[int][]PairwiseElimBallot)
	for i := range m.Electorate.Candidates {
		m.Buckets[i] = make([]PairwiseElimBallot, 0)
	}
}

//GetWinner returns the index of the winning candidate
func (m *PairwiseElimMethod) GetWinner() int {
	return m.Winner
}

//GetUtility returns the average utility for the winning candidate
func (m *PairwiseElimMethod) GetUtility() float64 {
	return m.Utility
}

//Run creates ballots and tabulates the winner
func (m *PairwiseElimMethod) Run() {

	//voters create ballots
	for i := range m.Electorate.Voters {
		if m.Electorate.Voters[i].Strategic {
			m.Ballots[i] = m.VoteStrategic(&m.Electorate.Voters[i])
		} else {
			m.Ballots[i] = m.Vote(&m.Electorate.Voters[i])
		}
	}

	var loser int
	var ballots []PairwiseElimBallot
	m.sortBallots(m.Ballots)
	//eliminate candidates (buckets) until only 1 remains
	for {
		if len(m.Buckets) == 1 {
			//there is a winner
			break
		}

		//identify loser, remove them, allocate their ballots
		loser = m.findLoser()
		ballots = m.Buckets[loser]
		delete(m.Buckets, loser)
		m.sortBallots(ballots)
	}

	//key of remaining candidate element is the winning candidate's index
	for k := range m.Buckets {
		m.Winner = k
	}

	m.calcUtility()
}

//returns index of next candidate to be eliminated
func (m *PairwiseElimMethod) findLoser() int {
	//see if there's a condorcet loser, return if there is
	isLoser, loser := m.Electorate.findGroupCondorcetLoser(m.remainingIndices())

	if isLoser {
		return loser
	}

	//if no condorcet loser, return whoever has the most hate ballots
	return m.findMostBallots()
}

//return index of candidate with most ballots
func (m *PairwiseElimMethod) findMostBallots() int {
	winner := -1
	most := 0

	for k := range m.Buckets {
		if len(m.Buckets[k]) > most {
			most = len(m.Buckets[k])
			winner = k
		}
	}

	return winner
}

//returns a slice with the indices of the remaining candidates
func (m *PairwiseElimMethod) remainingIndices() []int {
	candidates := make([]int, 0)
	for k := range m.Buckets {
		candidates = append(candidates, k)
	}

	return candidates
}

//move each ballot in slice provided to the bucket of the next remaining candidate
func (m *PairwiseElimMethod) sortBallots(ballots []PairwiseElimBallot) {
	//fmt.Println("sorting ballots")
	for i := range ballots {
		for {
			//decrement choice on ballot (move up from bottom)
			ballots[i].LastChoice--

			//see if that candidate remains
			_, ok := m.Buckets[ballots[i].LastChoice]

			//if that candidate's index still exists, add ballot to their slice of ballots
			if ok {
				m.Buckets[ballots[i].LastChoice] = append(m.Buckets[ballots[i].LastChoice], ballots[i])
				break

				//if not, check to see if this ballot is expired and if so, discard ballot
			} else if ballots[i].LastChoice < 0 {
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
func (m *PairwiseElimMethod) calcUtility() {
	u := 0.0
	for _, v := range m.Electorate.Voters {
		u += v.Utilities[m.Winner]
	}

	m.Utility = u / float64(len(m.Electorate.Voters))
}

//Vote creates a ballot for an honest voter
func (m *PairwiseElimMethod) Vote(v *Voter) PairwiseElimBallot {
	ballot := PairwiseElimBallot{Choices: make([]int, 0), LastChoice: len(m.Buckets)}

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
func (m *PairwiseElimMethod) VoteStrategic(v *Voter) PairwiseElimBallot {
	//strategic ranked voters will rank their preferred major candidate first and the other major candidate last
	preferredMajor := findFavoriteMajor(v.Utilities, m.Electorate.Candidates)
	otherMajor := findOtherMajor(preferredMajor, m.Electorate.Candidates)

	ballot := PairwiseElimBallot{Choices: make([]int, 0), LastChoice: len(m.Buckets)}

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
