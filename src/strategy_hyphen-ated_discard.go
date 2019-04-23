package main

func (d *Hyphenated) Discard(g *Game) *Action {
	// First, check for known-safe discards (from left-to-right)
	p := g.Players[g.ActivePlayer]
	for i := len(p.Hand) - 1; i >= 0; i-- {
		c := p.Hand[i]
		if d.Cards[c.Order].Trash {
			return &Action{
				Type:   actionTypeDiscard,
				Target: i,
			}
		}
	}

	chop := d.Players[g.ActivePlayer].Chop
	return &Action{
		Type:   actionTypeDiscard,
		Target: p.Hand[chop].Order,
	}
}
