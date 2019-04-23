package main

func (d *Hyphenated) CheckPlayable(g *Game) *Action {
	// Look through our whole hand and make a list of all the playable cards
	p := g.Players[d.Us]
	playables := make([]*Card, 0)
	for _, c := range p.Hand {
		if d.Cards[c.Order].Playable {
			playables = append(playables, c)
		}
	}
	if len(playables) == 0 {
		return nil
	}

	// Always play the left-most one
	// TODO Implement Priority
	cardToPlay := playables[0]
	return &Action{
		Type:   actionTypePlay,
		Target: cardToPlay.Order,
	}
}
