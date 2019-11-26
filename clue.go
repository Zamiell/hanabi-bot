package main

import (
	"strconv"
)

type Clue struct {
	Type  int `json:"type"` // 0 is a rank clue, 1 if a clue color
	Value int `json:"value"`
}

func (c *Clue) Name(g *Game) string {
	if c.Type == clueTypeRank {
		return strconv.Itoa(c.Value)
	} else if c.Type == clueTypeColor {
		return variants[g.Variant].ClueColors[c.Value]
	}

	return ""
}
