package main

func HyphenatedActionHappenedClue(d *Hyphenated, g *Game, a *Action) {
	p := g.Players[a.Target]

	for _, c := range p.Hand {
		hc := d.Cards[c.Order]

		if c.JustTouched &&
			c.Revealed &&
			(c.IsPlayable(g) || hc.IsDelayedPlayable(g, d)) {

			hc.Playable = true
		}
	}
}
