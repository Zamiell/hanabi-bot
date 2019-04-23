package main

func (d *Hyphenated) Discard(g *Game) *Action {
	// First, check for known-safe discards (from left-to-right)
	p := g.Players[d.Us]
	for i := len(p.Hand) - 1; i >= 0; i-- {
		c := p.Hand[i]
		if d.Cards[c.Order].Trash {
			return &Action{
				Type:   actionTypeDiscard,
				Target: i,
			}
		}
	}

	hp := d.Players[d.Us]
	return &Action{
		Type:   actionTypeDiscard,
		Target: hp.GetChop(g).Order,
	}
}
