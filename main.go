package main

import (
	"fmt"
	"sync"
	//"github.com/kr/pretty"
)

func main() {
	params := readParams()

	//create electorates
	electorates := createElectorates(params)

	//create methods for each electorate
	for i := range electorates {

		//create plurality method
		pm := PluralityMethod{}
		electorates[i].Methods["Plurality"] = &pm
		pm.Create(&electorates[i])

		//create approval method
		am := ApprovalMethod{}
		electorates[i].Methods["Approval"] = &am
		am.Create(&electorates[i])

		im := IRVMethod{}
		electorates[i].Methods["IRV"] = &im
		im.Create(&electorates[i])

	}

	var wg sync.WaitGroup

	//run methods in each electorate
	for i := range electorates {
		for name := range electorates[i].Methods {
			wg.Add(1)
			go electorates[i].Methods[name].Run(&wg)
		}
	}

	wg.Wait()

	efficiencies := make(map[string]float64)
	numEfficiencies := 0.0
	condorcets := make(map[string]float64)
	numCondorcets := 0.0

	for _, e := range electorates {
		//fmt.Println("-----")
		//printReport(e)

		r := e.GetReport()
		numEfficiencies += 1.0
		if e.CondorcetWinner > -1 {
			numCondorcets += 1.0
		}

		for n, l := range r.Lines {
			efficiencies[n] += l.Efficiency
			if e.CondorcetWinner > -1 {
				condorcets[n] += float64(l.Condorcet)
			}
		}
	}

	for n, eff := range efficiencies {
		eff = eff / numEfficiencies
		con := condorcets[n] / numCondorcets
		fmt.Printf("%s: %.3f %.2f \n", n, eff, con)
	}

}

func printReport(e Electorate) {
	r := e.GetReport()
	//fmt.Printf("%+v \n", r)
	fmt.Printf("Voters: %v\n", r.NumVoters)
	fmt.Printf("Candidates: %v\n", r.NumCandidates)
	fmt.Printf("Utility: %s\n", candidateInfo(r.UtilityWinner, e))
	fmt.Printf("Condorcet: %s\n", candidateInfo(r.CondorcetWinner, e))
	for name, l := range r.Lines {
		fmt.Printf("%s: %s, %.2f, %v \n", name, candidateInfo(l.Winner, e), l.Efficiency, l.Condorcet)
	}
}

func candidateInfo(i int, e Electorate) string {
	var major string
	var name string

	if i < 0 {
		major = ""
		name = "none"
	} else {
		name = e.Candidates[i].Name
		if e.Candidates[i].Major {
			major = "(major)"
		} else {
			major = ""
		}
	}

	return fmt.Sprintf("%s %s", name, major)
}
