package main

import (
	"strconv"
)

/*
	Main functions
*/

func HyphenatedActionAnnouncedClue(d *Hyphenated, g *Game, a *Action) {
	interpretation := d.GetClueInterpretation(g, a)
	clue := NewClue(a)
	if interpretation == hyphenClueTypePlay {
		focusedCard := d.GetClueFocus(g, a.Target, clue)
		if focusedCard == nil {
			logger.Fatal("The focused card of a play clue is nil.")
			return
		}
		hc := d.Cards[focusedCard.Order]
		hc.Playable = true
		logger.Debug(g.Players[d.Us].Name + " marked " + focusedCard.Name() + " on slot " + strconv.Itoa(focusedCard.Slot) + " as playable.")
	}

	//hp := d.Players[a.Target]
	//hp.UpdateHandGTP(g, d) // TODO this needs to be done to all players, not just the player receiving the clue
}

/*
	Subroutines
*/

func (d *Hyphenated) GetClueInterpretation(g *Game, a *Action) int {
	clue := NewClue(a)
	focusedCard := d.GetClueFocus(g, a.Target, clue)
	hp := d.Players[a.Target]

	if focusedCard == hp.GetChop(g, d) {
		// It is focused on the chop, so check to see if this could be a Save Clue.
		if a.Type == ActionTypeRankClue {
			// 2 Saves and 5 Saves.
			if a.Value == 2 || a.Value == 5 {
				return hyphenClueTypeSave
			}
		} else if a.Type == ActionTypeColorClue {
			// Check to see if the focused chop card "matches" anything in the discard pile.
			for _, c := range g.DiscardPile {
				if c.NeedsToBePlayed(g) &&
					focusedCard.CouldBeSuit(c.Suit) &&
					focusedCard.CouldBeRank(c.Rank) {

					return hyphenClueTypeSave
				}
			}
		}
	}
	return hyphenClueTypePlay
}
