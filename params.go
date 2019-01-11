package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

//AppParams holds all run parameters specified in params.json
type AppParams struct {
	NumElectorates     int      //the number of unique Electorates to generate and test
	MinVoters          int      //lower limit of randomly chosen size of electorate
	MaxVoters          int      //upper limit of randomly chosen size of electorate
	StrategicVoters    float64  //chance that a voter is "strategic"
	MinCandidates      int      //lower limit of randomly chosen number of candidates
	MaxCandidates      int      //upper limit of randomly chosen number of candidates
	NumMajorCandidates int      //the number of candidates representing "major parties". This should be either 0 or 2
	Methods            []string //list of methods that should be included in the analysis
	NumAxes            int      //the number of ideological axis that voters and candidates should align to
	Names              []string //list of all possible names for candidates. Must be at least as long as MaxCandidates
	NumWorkers         int      //number of concurrent workers to spawn for processing elections
}

func readParams() AppParams {
	//load public values stored in params.json
	raw, err := ioutil.ReadFile("params.json")
	if err != nil {
		panic(err)
	}

	var params AppParams
	err = json.Unmarshal(raw, &params)
	if err != nil {
		panic(err)
	}

	return params
}

func printParams(params *AppParams) {
	fmt.Println("Electorates:", params.NumElectorates)
	fmt.Println("Voters:", params.MinVoters, "to", params.MaxVoters)
	fmt.Println("Strategic Voters:", params.StrategicVoters)
	fmt.Println("Candidates:", params.MinCandidates, "to", params.MaxCandidates)
	fmt.Println("Methods:", params.Methods)
	fmt.Println("Axes:", params.NumAxes)
	fmt.Println(params.NumWorkers, "workers")
}

//returns true if the specified method is included in the app params
func includeMethod(method string, params *AppParams) bool {
	for _, m := range params.Methods {
		if m == method {
			return true
		}
	}

	return false
}
