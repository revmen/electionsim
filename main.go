package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
	//"github.com/kr/pretty"
)

func main() {
	start := time.Now()

	//get user values from params.json
	params := readParams()

	//random source that needs to be protected if used concurrently
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var mu sync.Mutex

	//create job channels and workers
	startChan := make(chan bool, params.NumWorkers)
	reviewChan := make(chan *Electorate, params.NumWorkers)
	summaryChan := make(chan string, params.NumWorkers)

	//start workers
	for i := 0; i < params.NumWorkers; i++ {
		go runWorker(&params, startChan, reviewChan, r, &mu)
	}

	go summaryWorker(&params, reviewChan, summaryChan)

	//starting the startWorker will begin the analysis
	go startWorker(&params, startChan)

	//wait for results, which will be printed until a value comes over the doneChan
	for line := range summaryChan {
		fmt.Println(line)
	}

	close(startChan)
	close(reviewChan)

	elapsed := time.Since(start)
	fmt.Printf("Analysis took %s \n", elapsed)
}

//promps runWorker to start jobs at a metered pace
func startWorker(params *AppParams, startChan chan bool) {
	//start a job for each electorate
	for i := 0; i < params.NumElectorates; i++ {
		startChan <- true
	}
}

//worker that pulls an electorate from the the queue and processes it
func runWorker(params *AppParams, startChan <-chan bool, reviewChan chan<- *Electorate, r *rand.Rand, mu *sync.Mutex) {

	for range startChan {
		//create electorate
		e := makeElectorate(params, r, mu)

		//create methods
		pm := PluralityMethod{}
		e.Methods["Plurality"] = &pm
		pm.Create(&e)

		am := ApprovalMethod{}
		e.Methods["Approval"] = &am
		am.Create(&e)

		im := IRVMethod{}
		e.Methods["IRV"] = &im
		im.Create(&e)

		//run methods
		for name := range e.Methods {
			e.Methods[name].Run()
		}

		reviewChan <- &e
	}
}

func summaryWorker(params *AppParams, reviewChan chan *Electorate, summaryChan chan string) {
	//create summary containers
	efficiencies := make(map[string]float64)
	numEfficiencies := 0.0
	condorcets := make(map[string]float64)
	numCondorcets := 0.0
	numCompleted := 0

	//extract results from completed electorates
	for e := range reviewChan {
		//add results to summaries
		r := e.GetReport()

		numEfficiencies += 1.0
		if e.CondorcetWinner > -1 {
			numCondorcets += 1.0
		}

		for m, l := range r.Lines {
			efficiencies[m] += l.Efficiency
			if e.CondorcetWinner > -1 {
				condorcets[m] += float64(l.Condorcet)
			}
		}

		numCompleted++

		if numCompleted >= params.NumElectorates {
			break
		}
	}

	//complete summary and pass text lines to main process
	for n, eff := range efficiencies {
		eff = eff / numEfficiencies
		con := condorcets[n] / numCondorcets
		summaryChan <- fmt.Sprintf("%s: %.3f %.2f", n, eff, con)
	}

	//signal completion of study by closing the summaryChan
	close(summaryChan)
}

func printReport(e *Electorate) {
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

func candidateInfo(i int, e *Electorate) string {
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
