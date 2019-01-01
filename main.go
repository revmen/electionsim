package main

import (
	"fmt"
	//"github.com/kr/pretty"
)

func main() {
	params := readParams()

	electorates := createElectorates(params)

	for i, e := range electorates {
		//create and run plurality method
		pm := PluralityMethod{}
		electorates[i].Methods["Plurality"] = &pm
		pm.Create(&e)
		pm.Run()

		//create and run approval method
		am := ApprovalMethod{}
		electorates[i].Methods["Approval"] = &am
		am.Create(&e)
		am.Run()
	}

	for _, e := range electorates {
		fmt.Println("-----")

		e.findUtilityWinner()
		fmt.Println("Utility", candidateInfo(e.UtilityWinner, e))

		e.findCondorcetWinner()
		fmt.Println("Condorcet", candidateInfo(e.CondorcetWinner, e))

		for name, method := range e.Methods {
			fmt.Println(name, candidateInfo(method.GetWinner(), e))

		}

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
