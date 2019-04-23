package main

import (
	"strconv"
)

func (d *Hyphenated) FindClueFocus(g *Game, i int, clue *Clue) *Card {
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

		for _, j := range []int{clueTypeRank, clueTypeColor} {
			if j == clueTypeRank {
				for _, k := range variants[g.Variant].Ranks {
					clue := d.CheckViableClue(g, i, j, k)
					if clue != nil {
						viableClues = append(viableClues, clue)
					}
				}
			} else if j == clueTypeColor {
				for _, k := range variants[g.Variant].Suits {
					clue := d.CheckViableClue(g, i, j, k.Index)
					if clue != nil {
						viableClues = append(viableClues, clue)
					}
				}
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
	c := d.FindClueFocus(g, i, clue)
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
