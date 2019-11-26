package main

// We store extra information about each player
type HyphenPlayer struct {
	Index int
}

func (hp *HyphenPlayer) GetChop(g *Game, d *Hyphenated) *Card {
	p := g.Players[hp.Index]
	for _, c := range p.Hand {
		hc := d.Cards[c.Order]

		if !c.Touched && !hc.ChopMoved {
			return c
		}
	}

	// The hand is fully-clued, so just make the chop the left-most card
	return p.Hand[len(p.Hand)-1]
}

// Now that a card has been clued,
// update the possibilities for each card in this player's hand based on Good Touch Principle
func (hp *HyphenPlayer) UpdateHandGTP(g *Game, d *Hyphenated) {
	p := g.Players[hp.Index]

	for _, c := range p.Hand {
		hc := d.Cards[c.Order]

		// TODO: don't mark a card as playable if it is known trash
		if c.Touched && !hc.Playable {
			isplayable := true
			for _, s := range c.PossibleSuits {
				for _, r := range c.PossibleRanks {
					if r > g.Stacks[s.Index]+1 {
						isplayable = false
						break
					}
				}
			}
			if isplayable {
				log.Info(c.Name() + " in " + p.Name + "'s hand is playable now from Good Touch Principle.")
				hc.Playable = true
			}
		}
	}
}
