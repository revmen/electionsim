package main

import (
	"fmt"
	//"github.com/kr/pretty"
)

func main() {
	params := readParams()

	electorates := createElectorates(params)

	for _, e := range electorates {
		e.UtilityWinner, e.MaxUtility = findUtilityWinner(e)
		fmt.Println("utility winner", e.UtilityWinner)
		e.CondorcetWinner = findCondorcetWinner(e)
		fmt.Println("condorcet winner", e.CondorcetWinner)
		// if "no condorcet winner" == e.CondorcetWinner {
		// 	//fmt.Println(i, len(e.Candidates), "candidates")
		// 	//fmt.Println("utilty", uName)
		// 	//fmt.Println("condorcet", cName)
		// 	//fmt.Printf("%# v", pretty.Formatter(e))
		// 	printPluralityVotes(e)
		// }
		pe := PluralityElection{}
		pe.DoElection(e)
		fmt.Println("plurality winner", pe.Winner)
	}

}
