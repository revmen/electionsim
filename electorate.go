package main

import (
	"math/rand"
	"sync"
)

//Electorate is a collection of Voters and Candidates
type Electorate struct {
	Voters          []Voter           //slice of all voters in electorate
	Candidates      []Candidate       //slice of all candidates
	MaxUtility      float64           //average utility per voter for max utility candidate
	UtilityWinner   int               //index of max utility candidate
	CondorcetWinner int               //index of the condorcet winner
	Methods         map[string]Method //map of Method interfaces with name of election method as key
}

//Voter represents an individual voter with unique alignments in each axis and a flag for whether the voter is "strategic"
type Voter struct {
	Alignments        []float64 //the ideological alignment of the voter based on scores in axes
	Strategic         bool      //whether or not hte voter votes "strategically"
	Utilities         []float64 //the utilty the voter has for each candidate
	ApprovalThreshold float64   //the utility threshold required for a voter to be OK with a candidate
}

//Candidate is a single ballot choice with specific alignments
type Candidate struct {
	Name       string
	Alignments []float64
	Major      bool //true if the candidate is from a "major party"
}

//Report is a summary of the performance of all methods run in an electorate
type Report struct {
	NumVoters       int                   //number of voters in the electorate
	NumCandidates   int                   //number of candidates on the ballot
	CondorcetWinner int                   //index of condorcet winner
	UtilityWinner   int                   //index of highest utility candidate
	Lines           map[string]ReportLine //summary for each method, name of method as key
}

//ReportLine is a single line in a report, covering one voting method
type ReportLine struct {
	Winner     int     //index of the winning Candidate
	Efficiency float64 //the fraction of maximum possible efficiency achieved with the winning candidate
	Condorcet  int     //whether the Condorcet winner was elected. 0 for false, 1 for true, -1 means there was no Condorcet winner.
}

//GetReport creates and returns a Report, which is a summary of the performance of methods tested for this electorate
func (e *Electorate) GetReport() Report {
	r := Report{
		NumVoters:       len(e.Voters),
		NumCandidates:   len(e.Candidates),
		CondorcetWinner: e.CondorcetWinner,
		UtilityWinner:   e.UtilityWinner,
		Lines:           make(map[string]ReportLine),
	}

	for name, m := range e.Methods {
		c := -1

		//mark whether the condorcet winner was matched by this method
		//value of -1 means there is no condorcet winner
		if e.CondorcetWinner > -1 {
			if e.CondorcetWinner == m.GetWinner() {
				c = 1
			} else {
				c = 0
			}
		}

		//add the method's result to the report
		r.Lines[name] = ReportLine{
			Winner:     m.GetWinner(),
			Efficiency: m.GetUtility() / e.MaxUtility,
			Condorcet:  c,
		}
	}

	return r
}

func makeElectorate(params *AppParams, r *rand.Rand, mu *sync.Mutex) Electorate {
	e := Electorate{}

	//lock random number generator, decide number of candidates, decide number of voters
	mu.Lock()
	numCandidates := r.Intn(params.MaxCandidates-params.MinCandidates+1) + params.MinCandidates
	numVoters := r.Intn(params.MaxVoters-params.MinVoters+1) + params.MinVoters
	mu.Unlock()

	//create candidates
	e.Candidates = make([]Candidate, numCandidates)
	for i := 0; i < numCandidates; i++ {
		if i < params.NumMajorCandidates {
			e.Candidates[i] = makeMajorCandidate(params.Names[i], params.NumAxes, i, r, mu)
		} else {
			e.Candidates[i] = makeCandidate(params.Names[i], params.NumAxes, r, mu)
		}
	}

	//create Voters
	e.Voters = make([]Voter, numVoters)
	for i := 0; i < numVoters; i++ {
		e.Voters[i] = makeVoter(params.NumAxes, params.StrategicVoters, e.Candidates, r, mu)
	}

	//create map for methods
	e.Methods = make(map[string]Method)

	//determine the utility and condorcet winners for this electorate
	e.findUtilityWinner()
	e.findCondorcetWinner()

	return e
}

//create a single voter
func makeVoter(numAxes int, strategicChance float64, candidates []Candidate, r *rand.Rand, mu *sync.Mutex) Voter {
	//create the ideological axes
	axes := make([]float64, numAxes)

	//lock the random number generator, populate the axes, decide whether voter is strategic
	mu.Lock()
	for i := 0; i < len(axes); i++ {
		axes[i] = r.Float64()
	}
	isStrategic := r.Float64() <= strategicChance
	mu.Unlock()

	//create slice of utilities used to hold voter's utility from each candidate
	utilities := make([]float64, len(candidates))

	//assemble Voter struct
	v := Voter{
		Alignments:        axes,
		Strategic:         isStrategic,
		Utilities:         utilities,
		ApprovalThreshold: 0.5,
	}

	//determine voter's utilities
	for i, c := range candidates {
		utilities[i] = utility(v, c)
	}

	return v
}

//creates a single candidate that is not a "major"
func makeCandidate(name string, numAxes int, r *rand.Rand, mu *sync.Mutex) Candidate {
	//create the ideological axes
	axes := make([]float64, numAxes)

	//lock the random number generator and populate the axes
	mu.Lock()
	for i := 0; i < len(axes); i++ {
		axes[i] = r.Float64()
	}
	mu.Unlock()

	//populate and return Candidate struct
	c := Candidate{
		Alignments: axes,
		Name:       name,
		Major:      false,
	}

	return c
}

//creates a single candidate that is a "major"
func makeMajorCandidate(name string, numAxes int, index int, r *rand.Rand, mu *sync.Mutex) Candidate {
	//create the ideological axes
	axes := make([]float64, numAxes)

	//which quadrant/octant candidate is in based on whether index is even or odd
	zone := index % 2
	min := float64(zone) * 0.5
	max := float64(zone)*0.5 + 0.5

	//lock the random number generator and populate the axes
	//a major candidate has all of their alignments in the same quadrant/octant, where axis crossing are at 0.5
	mu.Lock()
	for i := 0; i < len(axes); i++ {
		axes[i] = min + r.Float64()*(max-min)
	}
	mu.Unlock()

	//populate and return Candidate struct
	c := Candidate{
		Alignments: axes,
		Name:       name,
		Major:      true,
	}

	return c
}

//AverageUtilityOf returns the average utility for the candidate at the specified index
func (e *Electorate) AverageUtilityOf(candidateIndex int) float64 {
	sum := 0.0
	for _, voter := range e.Voters {
		sum += voter.Utilities[candidateIndex]
	}

	return sum / float64(len(e.Voters))
}
