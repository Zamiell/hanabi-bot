package main

import (
	"strconv"
)

type Card struct {
	Suit          *Suit
	Rank          int
	Order         int
	Clues         []*CardClue
	PossibleSuits []*Suit
	PossibleRanks []int
	PossibleCards map[string]int
	Revealed      bool
}
type CardClue struct {
	Type     int
	Value    int
	Positive bool
}

func (c *Card) Name(g *Game) string {
	name := c.Suit.Name + " " + strconv.Itoa(c.Rank)
	return name
}

func (c *Card) RemovePossibility(suit *Suit, rank int, removeAll bool) {
	// Every card has a possibility map that maps card identities to count
	mapIndex := suit.Name + strconv.Itoa(rank)
	cardsLeft := c.PossibleCards[mapIndex]
	if cardsLeft > 0 {
		// Remove one or all possibilities for this card,
		// (depending on whether the card was clued
		// or if we saw someone draw a copy of this card)
		cardsLeft := cardsLeft - 1
		if removeAll {
			cardsLeft = 0
		}
		c.PossibleCards[mapIndex] = cardsLeft
	}
}
