package main

import (
	"strconv"
)

type Clue struct {
	Type  int `json:"type"` // 0 if a color clue, 1 if a rank color
	Value int `json:"value"`
}

func NewClue(a *Action) *Clue {
	return &Clue{
		// A color clue is action type 2
		// A rank clue is action type 3
		// Remap these to 0 and 1, respectively
		Type:  a.Type - 2,
		Value: a.Value,
	}
}

func (c *Clue) Name(g *Game) string {
	if c.Type == ClueTypeRank {
		return strconv.Itoa(c.Value)
	} else if c.Type == ClueTypeColor {
		return variants[g.Variant].ClueColors[c.Value]
	}

	return ""
}
