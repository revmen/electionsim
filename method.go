package main

//Method interface to allow general functions to call functions in any election mehtod
type Method interface {
	Create(*Electorate)
	Run()
	GetWinner() int
}
