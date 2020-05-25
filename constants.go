package main

const (
	ActionTypePlay = iota
	ActionTypeDiscard
	ActionTypeColorClue
	ActionTypeRankClue
)

const (
	ClueTypeColor = iota
	ClueTypeRank
)

const (
	EndConditionInProgress = iota
	EndConditionNormal
	EndConditionStrikeout
)

const (
	// The maximum amount of clues (and the amount of clues that players start a game with)
	MaxClueNum = 8

	// The maximum amount of strikes/misplays allowed before the game ends
	MaxStrikeNum = 3
)
