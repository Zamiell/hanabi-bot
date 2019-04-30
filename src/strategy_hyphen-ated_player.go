package main

// We store extra information about each player
type HyphenPlayer struct {
	Index int
	Chop  int // Equal to the index in the hand of the chop card
}

func (hp *HyphenPlayer) GetChop(g *Game) *Card {
	p := g.Players[hp.Index]
	return p.Hand[hp.Chop]
}

func (hp *HyphenPlayer) UpdateChop(g *Game) {
	p := g.Players[hp.Index]
	for i, c := range p.Hand {
		if !c.IsClued() {
			hp.Chop = i
		}
	}

	// The hand is fully-clued, so just make the chop the left-most card
	hp.Chop = len(p.Hand) - 1
}
