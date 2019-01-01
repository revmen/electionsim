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
	UtilityWinner   int         //index of max utility candidate
	CondorcetWinner int         //index of the condorcet winner
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
	Major      bool //true if the candidate is from a "major party"
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
		if i < params.NumMajorCandidates {
			candidates[i] = makeMajorCandidate(params.Names[i], params.NumAxes, i, r)
		} else {
			candidates[i] = makeCandidate(params.Names[i], params.NumAxes, r)
		}
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
		Major:      false,
	}

	return c
}

func makeMajorCandidate(name string, numAxes int, index int, r *rand.Rand) Candidate {
	axes := make([]float64, numAxes)

	//which quadrant/octant candidate is in based on whether index is even or odd
	zone := index % 2
	min := float64(zone) * 0.5
	max := float64(zone)*0.5 + 0.5

	//a major candidate has all of their alignments in the same quadrant/octant, where axis crossing are at 0.5
	for i := 0; i < len(axes); i++ {
		axes[i] = min + r.Float64()*(max-min)
	}

	c := Candidate{
		Alignments: axes,
		Name:       name,
		Major:      true,
	}

	return c
}
