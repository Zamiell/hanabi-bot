package main

func (d *Hyphenated) CheckPlayClues(g *Game) *Action {
	playableCardsToGet := make([]*Card, 0)
	for i, p := range g.Players {
		if i == d.Us {
			continue
		}

		for _, c := range p.Hand {
			hc := d.Cards[c.Order]
			if (c.IsPlayable(g) || hc.IsDelayedPlayable(g, d)) &&
				// This card already has a clue on it,
				// so cluing it again would not satisfy Minimum Clue Value Principle
				!c.Touched &&
				// This card was already gotten via a Finesse
				!hc.Playable {

				playableCardsToGet = append(playableCardsToGet, c)
			}
		}
	}

	viableClues := make([]*PossibleClue, 0)
	for _, c := range playableCardsToGet {
		var clue *PossibleClue
		var alreadyClued bool

		// Rank clue
		alreadyClued = false
		for _, clue := range viableClues {
			if clue.Clue.Type == clueTypeRank &&
				clue.Clue.Value == c.Rank &&
				clue.Target == c.Holder {

				alreadyClued = true
				break
			}
		}
		if !alreadyClued {
			clue = d.CheckViableClue(g, c.Holder, clueTypeRank, c.Rank, hyphenClueTypePlay)
			if clue != nil {
				viableClues = append(viableClues, clue)
			}
		}

		// Color clue
		alreadyClued = false
		for _, clue := range viableClues {
			if clue.Clue.Type == clueTypeColor &&
				clue.Clue.Value == c.Suit.ColorValue(g) &&
				clue.Target == c.Holder {

				alreadyClued = true
				break
			}
		}
		if !alreadyClued {
			clue = d.CheckViableClue(g, c.Holder, clueTypeColor, c.Suit.ColorValue(g), hyphenClueTypePlay)
			if clue != nil {
				viableClues = append(viableClues, clue)
			}
		}
	}

	// Prefer the clues that touch the greatest amount of cards
	for i := g.GetHandSize(); i >= 0; i-- {
		for _, viableClue := range viableClues {
			if viableClue.CardsClued == i {
				return &Action{
					Type:   actionTypeClue,
					Clue:   viableClue.Clue,
					Target: viableClue.Target,
				}
			}
		}
	}

	return nil
}
