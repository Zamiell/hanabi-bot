package main

const (
	actionTypeClue = iota
	actionTypePlay
	actionTypeDiscard
)

const (
	clueTypeRank = iota
	clueTypeColor
)

const (
	endConditionInProgress = iota
	endConditionNormal
	endConditionStrikeout
)

const (
	// The maximum amount of clues (and the amount of clues that players start a game with)
	maxClues = 8
)
