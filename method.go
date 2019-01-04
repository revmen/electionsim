package main

import (
	"sync"
)

//Method interface to allow general functions to call functions in any election mehtod
type Method interface {
	Create(*Electorate)
	Run(*sync.WaitGroup)
	GetWinner() int
	GetUtility() float64
}
