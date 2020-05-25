package main

type CardClue struct {
	Type     int
	Value    int
	Positive bool
}

func (c *CardClue) Name(g *Game) string {
	name := ""
	if c.Type == ClueTypeRank {
		name = string(c.Value)
	} else if c.Type == ClueTypeColor {
		name = variants[g.Variant].ClueColors[c.Value]
	}
	if !c.Positive {
		name = "-" + name
	}
	return name
}
