package main

// Add the 3 functions to the function maps
func dumbInit() {
	name := "Dumb"
	stratStart[name] = dumbStart
	stratGetAction[name] = dumbGetAction
	stratActionHappened[name] = dumbActionHappened
}

// Called at the beginning of the game
func dumbStart(g *Game) {
}

// Called when it gets to the player's turn
// Should returns the action that the player should perform
func dumbGetAction(g *Game) *Action {
	// Get the our slot 1 card
	hand := g.Players[g.ActivePlayer].Hand
	card := hand[len(hand)-1]

	return &Action{
		Type:   actionTypePlay,
		Target: card.Order,
	}
}

// Called when an action happens
func dumbActionHappened(g *Game) {
}
