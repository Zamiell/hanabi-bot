package main

func (d *Hyphenated) CheckPlayableCards(g *Game) *Action {
	// Look through our whole hand and make a list of all the playable cards.
	p := g.Players[d.Us]
	playables := make([]*Card, 0)
	for _, c := range p.Hand {
		hc := d.Cards[c.Order]

		if hc.Playable {
			playables = append(playables, c)
		}
	}
	if len(playables) == 0 {
		return nil
	}

	// Always play the left-most one.
	// TODO Implement Priority
	cardToPlay := playables[0]
	return &Action{
		Type:   ActionTypePlay,
		Target: cardToPlay.Order,
	}
}
