package main

import ()

//Method interface to allow general functions to call functions in any election method
type Method interface {
	Create(*Electorate)
	Run()
	GetWinner() int
	GetUtility() float64
}
