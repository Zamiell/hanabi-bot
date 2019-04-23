package main

import (
	"strconv"
)

func (d *Hyphenated) GetClueFocus(g *Game, i int, clue *Clue) *Card {
	// First, look for freshly touched cards
	p := g.Players[i]
	freshCards := p.GetFreshCardsTouchedByClue(g, clue)
	if len(freshCards) == 1 {
		// The focus of the clue is on the only brand new card introduced
		return freshCards[0]
	}
	if len(freshCards) > 1 {
		// If one of the brand new cards introduced is on the chop, the focus is the chop
		for _, c := range freshCards { // TODO: Chop Moving will break this logic
			chop := d.Players[i].Chop
			if c == p.Hand[chop] {
				return c
			}
		}

		// Otherwise, the focus is the left-most of the freshly touched cards
		for _, c := range p.Hand {
			for _, ft := range freshCards {
				if ft == c {
					return c
				}
			}
		}
	}

	// If no brand new cards were introduced, the focus of the clue is the left-most card
	for _, c := range p.Hand {
		if c.JustTouched {
			return c
		}
	}

	return nil
}

func (d *Hyphenated) GetClueInterpretation(g *Game, a *Action) int {
	focusedCard := d.GetClueFocus(g, a.Target, a.Clue)

	hp := d.Players[a.Target]
	if focusedCard == hp.GetChop(g) {
		// It is focused on the chop, so check to see if this could be a Save Clue
		if a.Clue.Type == clueTypeRank {
			// 2 Saves and 5 Saves
			if a.Clue.Value == 2 || a.Clue.Value == 5 {
				return hyphenClueTypeSave
			}
		} else if a.Clue.Type == clueTypeColor {
			// Check to see if the focused chop card "matches" anything in the discard pile
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

func (d *Hyphenated) UpdateChop(g *Game, a *Action) {
	p := g.Players[a.Target]
	for i, c := range p.Hand {
		hasPositiveClue := false
		for _, clue := range c.Clues {
			if clue.Positive {
				hasPositiveClue = true
				break
			}
		}
		if !hasPositiveClue {
			d.Players[a.Target].Chop = i
		}
	}

	// The hand is fully-clued, so just make the chop the left-most card
	d.Players[a.Target].Chop = len(p.Hand) - 1
}

type PossibleClue struct {
	Clue       *Clue
	Target     int
	CardsClued int
}

func (d *Hyphenated) CheckPlayClues(g *Game) *Action {
	viableClues := make([]*PossibleClue, 0)
	for i := range g.Players {
		if i == d.Us {
			continue
		}
		// Rank clues
		for _, k := range variants[g.Variant].ClueRanks {
			clue := d.CheckViableClue(g, i, clueTypeRank, k)
			if clue != nil {
				viableClues = append(viableClues, clue)
			}
		}
		// Color clues
		for j, _ := range variants[g.Variant].ClueColors {
			clue := d.CheckViableClue(g, i, clueTypeColor, j)
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

func (d *Hyphenated) CheckViableClue(g *Game, i int, j int, k int) *PossibleClue {
	clue := &Clue{
		Type:  j,
		Value: k,
	}
	p := g.Players[i]
	touchedCards := p.GetCardsTouchedByClue(g, clue)

	// We are not allowed to give a clue that touches 0 cards in the hand
	if len(touchedCards) == 0 {
		return nil
	}

	// Check if any of the touched cards are duplicates of one another
	if len(touchedCards) >= 2 {
		for _, c := range touchedCards {
			for _, c2 := range touchedCards {
				if c == c2 {
					continue
				}
				if c.Suit == c2.Suit && c.Rank == c2.Rank {
					return nil
				}
			}
		}
	}

	// Check to see if the card would misplay if we clued it
	c := d.GetClueFocus(g, i, clue)
	if !c.IsPlayable(g) {
		return nil
	}

	// Check for Good Touch Principle
	freshCards := p.GetFreshCardsTouchedByClue(g, clue)
	for _, c := range freshCards {
		for i, p := range g.Players {
			for _, c2 := range p.Hand {
				if i == d.Us {
					// GTP 1/2 - Don't potentialy duplicate clued cards in our hand
					mapIndex := c.Suit.Name + strconv.Itoa(c.Rank)
					if c2.IsClued() && c2.PossibleCards[mapIndex] > 0 {
						return nil
					}
				} else {
					// GTP 2/2 - Don't duplicate cards in other players hands
					if c2.IsClued() && c.Suit == c2.Suit && c.Rank == c2.Rank {
						return nil
					}
				}
			}
		}
	}

	return &PossibleClue{
		Clue: &Clue{
			Type:  j,
			Value: k,
		},
		Target:     i,
		CardsClued: len(freshCards),
	}
}
