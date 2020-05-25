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
				// so cluing it again would not satisfy Minimum Clue Value Principle.
				!c.Touched &&
				// This card was already gotten via a Finesse.
				!hc.Playable {

				playableCardsToGet = append(playableCardsToGet, c)

				if c.IsPlayable(g) {
					logger.Debug(c.Name() + " in " + g.Players[c.Holder].Name + "'s " +
						"hand is playable...")
				} else if hc.IsDelayedPlayable(g, d) {
					logger.Debug(c.Name() + " in " + g.Players[c.Holder].Name + "'s " +
						"hand is delayed playable...")
				}
			}
		}
	}

	viableClues := make([]*PossibleClue, 0)
	for _, c := range playableCardsToGet {
		var clue *PossibleClue
		var alreadyClued bool

		// Color clue
		alreadyClued = false
		for _, clue := range viableClues {
			if clue.Clue.Type == ClueTypeColor &&
				clue.Clue.Value == c.Suit.ColorValue(g) &&
				clue.Target == c.Holder {

				alreadyClued = true
				break
			}
		}
		if !alreadyClued {
			clue = d.CheckViableClue(
				g,
				c.Holder,
				ClueTypeColor,
				c.Suit.ColorValue(g),
				hyphenClueTypePlay,
			)
			if clue != nil {
				viableClues = append(viableClues, clue)
			}
		}

		// Rank clue
		alreadyClued = false
		for _, clue := range viableClues {
			if clue.Clue.Type == ClueTypeRank &&
				clue.Clue.Value == c.Rank &&
				clue.Target == c.Holder {

				alreadyClued = true
				break
			}
		}
		if !alreadyClued {
			clue = d.CheckViableClue(
				g,
				c.Holder,
				ClueTypeRank,
				c.Rank,
				hyphenClueTypePlay,
			)
			if clue != nil {
				viableClues = append(viableClues, clue)
			}
		}
	}

	if len(viableClues) == 0 {
		return nil
	}

	// Prefer the clues that touch the greatest amount of cards.
	for i := g.GetHandSize(); i >= 0; i-- {
		for _, viableClue := range viableClues {
			if viableClue.CardsClued == i {
				var actionType int
				if viableClue.Clue.Type == ClueTypeColor {
					actionType = ActionTypeColorClue
				} else if viableClue.Clue.Type == ClueTypeRank {
					actionType = ActionTypeRankClue
				}
				return &Action{
					Type:   actionType,
					Target: viableClue.Target,
					Value:  viableClue.Clue.Value,
				}
			}
		}
	}

	return nil
}
