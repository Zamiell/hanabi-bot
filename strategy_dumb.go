package main

func dumbInit() {
	name := "Dumb"
	stratStart[name] = dumbStart
	stratGetAction[name] = dumbGetAction
	stratActionHappened[name] = dumbActionHappened
}

func dumbStart(g *Game) {
	// Called at the beginning of the game
}

func dumbGetAction(g *Game) *Action {
	// Called when it gets to the player's turn
	// Should returns the action that the player should perform
	return nil
}

func dumbActionHappened(g *Game) {
	// Called when an action happens
}
