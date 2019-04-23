package main

import (
	"strconv"
)

type Card struct {
	Suit  int
	Rank  int
	Order int
	Clues []CardClue
}

type CardClue struct {
	Type     int
	Value    int
	Positive bool
}

func (c *Card) Name(g *Game) string {
	name := variants[g.Variant].Suits[c.Suit].Name // The name of the suit that this card is
	name += " "
	name += strconv.Itoa(c.Rank)
	return name
}
