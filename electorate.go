package main

import (
	"math/rand"
	"time"
)

//Electorate is a collection of Voters and Candidates
type Electorate struct {
	Voters          []Voter     //slice of all voters in electorate
	Candidates      []Candidate //slice of all candidates
	MaxUtility      float64     //average utility per voter for max utility candidate
	UtilityWinner   string      //name of max utility candidate
	CondorcetWinner string      //name of the condorcet winner
}

//Voter represents an individual voter with unique alignments in each axis and a flag for whether the voter is "strategic"
type Voter struct {
	Alignments []float64 //the ideological alignment of the voter based on scores in axes
	Strategic  bool      //whether or not hte voter votes "strategically"
	Utilities  []float64 //the utilty the voter has for each candidate
}

//Candidate is a single ballot choice with specific alignments
type Candidate struct {
	Name       string
	Alignments []float64
}

func createElectorates(params AppParams) []Electorate {
	electorates := make([]Electorate, params.NumElectorates)

	for i := 0; i < params.NumElectorates; i++ {
		electorates[i] = makeElectorate(params)
	}

	return electorates
}

func makeElectorate(params AppParams) Electorate {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	numCandidates := r.Intn(params.MaxCandidates-params.MinCandidates+1) + params.MinCandidates
	candidates := make([]Candidate, numCandidates)
	for i := 0; i < numCandidates; i++ {
		candidates[i] = makeCandidate(params.Names[i], params.NumAxes, r)
	}

	numVoters := r.Intn(params.MaxVoters-params.MinVoters+1) + params.MinVoters
	voters := make([]Voter, numVoters)
	for i := 0; i < numVoters; i++ {
		voters[i] = makeVoter(params.NumAxes, params.StrategicVoters, candidates, r)
	}

	e := Electorate{
		Voters:     voters,
		Candidates: candidates,
	}

	return e
}

func makeVoter(numAxes int, strategicChance float64, candidates []Candidate, r *rand.Rand) Voter {
	axes := make([]float64, numAxes)

	for i := 0; i < len(axes); i++ {
		axes[i] = r.Float64()
	}

	utilities := make([]float64, len(candidates))

	v := Voter{
		Alignments: axes,
		Strategic:  isStrategic(strategicChance, r),
		Utilities:  utilities,
	}

	for i, c := range candidates {
		utilities[i] = utility(v, c)
	}

	return v
}

func isStrategic(strategicChance float64, r *rand.Rand) bool {
	v := r.Float64()
	return v <= strategicChance
}

func makeCandidate(name string, numAxes int, r *rand.Rand) Candidate {
	axes := make([]float64, numAxes)

	for i := 0; i < len(axes); i++ {
		axes[i] = r.Float64()
	}

	c := Candidate{
		Alignments: axes,
		Name:       name,
	}

	return c
}
