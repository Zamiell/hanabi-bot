package main

type HyphenPlayer struct {
	Index int
	Chop  int // Equal to the index in the hand of the chop card
}

func (hp *HyphenPlayer) GetChop(g *Game) *Card {
	p := g.Players[hp.Index]
	return p.Hand[hp.Chop]
}
