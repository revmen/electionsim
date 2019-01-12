package main

//Method interface to allow general functions to call functions in any election method
type Method interface {
	Create(*Electorate)
	Run()
	GetWinner() int
	GetUtility() float64
}

// SimpleMethod is a simpler method interface. It can be adapted into the normal Method using AdaptSimpleMethod().
type SimpleMethod interface {
	FindWinner(*Electorate) int
}

// AdaptedMethod adapts a SimpleMethod to work with the normal Method interface
type AdaptedMethod struct {
	internalMethod SimpleMethod
	electorate     *Electorate
	winner         int
	utility        float64
}

// AdaptSimpleMethod adapts a SimpleMethod to be compatible with the normal Method interface
func AdaptSimpleMethod(method SimpleMethod) AdaptedMethod {
	return AdaptedMethod{internalMethod: method}
}

// Create initializes the adapter
func (method *AdaptedMethod) Create(electorate *Electorate) {
	method.electorate = electorate
}

// Run calculates and records the index of the winner and their average utility
func (method *AdaptedMethod) Run() {
	method.winner = method.internalMethod.FindWinner(method.electorate)
	method.utility = method.electorate.AverageUtilityOf(method.winner)
}

// GetWinner is an accessor for the index of the winner found by Run()
func (method *AdaptedMethod) GetWinner() int {
	return method.winner
}

// GetUtility is an accessor for the average utility of the winner
func (method *AdaptedMethod) GetUtility() float64 {
	return method.utility
}
