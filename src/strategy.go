package main

var (
	stratStart          map[string]func(*Game)
	stratGetAction      map[string]func(*Game) *Action
	stratActionHappened map[string]func(*Game)
)

func stratInit() {
	stratStart = make(map[string]func(*Game))
	stratGetAction = make(map[string]func(*Game) *Action)
	stratActionHappened = make(map[string]func(*Game))

	// New strategies that are added need to add their initialization function here
	dumbInit()
}
