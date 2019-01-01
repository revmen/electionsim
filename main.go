package main

import (
	"fmt"
	//"github.com/kr/pretty"
)

func main() {
	params := readParams()

	electorates := createElectorates(params)

	for _, e := range electorates {
		fmt.Println("-----")
		e.UtilityWinner, e.MaxUtility = findUtilityWinner(e)
		fmt.Println("utility winner", candidateInfo(e.UtilityWinner, e))
		e.CondorcetWinner = findCondorcetWinner(e)
		fmt.Println("condorcet winner", candidateInfo(e.CondorcetWinner, e))
		// if "no condorcet winner" == e.CondorcetWinner {
		// 	//fmt.Println(i, len(e.Candidates), "candidates")
		// 	//fmt.Println("utilty", uName)
		// 	//fmt.Println("condorcet", cName)
		// 	//fmt.Printf("%# v", pretty.Formatter(e))
		// 	printPluralityVotes(e)
		// }
		pe := PluralityElection{}
		pe.DoElection(e)
		fmt.Println("plurality winner", candidateInfo(pe.Winner, e))
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
