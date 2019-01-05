# Rev's Election Simulator

Here is yet another election simulator. This is similar to and in the same spirit as Warren Smith's work. I wrote this as an exercise and to increase my own understanding of voting methods.

I am very interested in feedback, both from voting enthusiasts and from programmers. I'm an amateur in both ways and it's a guarantee that I've made mistakes.

This tool is written in Go, which is an easy-to-read language that has very useful concurrency tools. Go can take advantage of today's multi-thread processors and operating systems to crank through large tests faster.

The code is organized into files that each contain one major building block.

## Electorate
An Electorate is a single struct that includes all of the Voters and Candidates that would take part in an election. There are a few additional tools and values used for analysis. These structs and functions are found on electorate.go.

Voters and Candidates each have a unique set of randomly generated, ideological alignments across a user-defined number of axes. Each value is a float64 between 0 and 1.

The difference in geometric space between a Voter and a Candidate determine's the Voter's utility if the Candidate is elected. Utilities are normalized to 0 to 1.

Candidates can be normal or "major". Major candidates are used in the strategies of strategic voters. There can be either 0 or 2 major candidates. If there are major candidates, they will be created with alignments that fall in opposing quadrants, octants, etc. This means that one major candidate will have all of their alignment values greater than 0.5, and one will have all of their values below 0.5.

Voters are "honest" or "strategic." An honest voter votes based only on their preference. A strategic voter's rules involve the major candidates. The approximate fraction of strategic voters in an electorate can be set by the user.

If the number of major candidates is set to 0, the fraction of strategic voters should also be set to 0.

## Criteria
Currently, 2 criteria for success are considered: Utility Efficiency and Condorcet. Functions related to these are found in utility.go and condorcet.go.

Utility Efficiency is really the same thing as Bayesian Regret used in other simulators. The winning Candidate is compared to the Candidate that would have produced the highest overall utility. The total achieved utility across all voters is divided by the total possible utility. In many elections, the winning candidate and the "best" candidate will be the same, which means a Utility Efficiecny of 1.0. In some cases, the "best" candidate will not win, which will result in a lower efficiency. Over many simulations, an average efficiency can be calculated.

The Condorcet Winner is the candidate that wins every individual head-to-head matchup. There isn't always a Condorcet Winner. Over many simuluations, a likelihood of electing the Condorcet winner can be calculated.

In future versions I'd like to consider other, more complicated criteria. I'd also like to look for failures like non-monotonicity.

## Methods
On method.go is the Method interface. This interface allows the processing functions on main.go to call functions attached to any of the voting methods currently found on approval.go, plurality.go, and irv.go. Each Method has unique routines for creating and counting ballots.

Each Method has the logic needed to conduct its style of election in an Electorate. Each Method will have its own file and will have certain key functions in common with all other Methods.

Each Method is run for every Electorate.

To encourage others to submit new Methods, I've kept all of the complicated concurrency stuff in main.go and electorate.go. If you'd like to submit a Method, you should be able to copy any of the existing Methods and modify them appropriately.

## Parameters
The user-definable parameters are found in params.go and can be set by the user in params.json. If you're not planning on doing any programming and just want to run the software, params.json is the only file you should edit.

#### NumElectorates
The number of unique electorates to test each method against. A high number produces a more statistically significant result. In early testing, it seems that anything above 2000 doesn't change the results.

#### MinVoters and MaxVoters
The number of voters in an electorate will be randomly selected to be between these values. With a range of 10,000 to 30,000, early testing shows results that aren't really any different from higher values.

#### StrategicVoters
The chance that a voter will be "strategic". This should be a fraction between 0.0 and 1.0.

If there are no Major Candidates, this value MUST be 0.0 to produce meaningful results.

#### MinCandidates and MaxCandidates
These values set the range for the possible number of candidates for each electorate. Since we're comparing multi-candidate voting system, the Min value should be at least 3.

#### NumMajorCandidates
This tells the simulator how many major candidates should be created. This value should be either 0 or 2.

If this value is 0, then StrategicVoters MUST be 0.0.

#### NumAxes
The number of ideological axes on which each voter and candidates alignment will fall. Values of 2 or 3 provide plenty of room for there to be meaningful differentiation.

#### Names
Just a list of names for candidates to be used when observing results from individual elections. This type of analysis isn't currently included, so these names are mostly for testing.

#### NumWorkers
This sets the number of concurrent workers that will be used to process elections concurrently. The higher this number, the more processor memory the program will use. This number should be lower for older or simpler machines.

