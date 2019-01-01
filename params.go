package main

import (
	"encoding/json"
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
	NumMajorCandidates int      //the number of candidates representing "major parties"
	NumAxes            int      //the number of ideological axis that voters and candidates should align to
	Names              []string //list of all possible names for candidates. Must be at least as long as MaxCandidates
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
