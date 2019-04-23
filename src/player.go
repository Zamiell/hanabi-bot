package main

import (
	"strconv"
)

type Player struct {
	Name     string
	Hand     []*Card
	Notes    []string
	Strategy *Strategy
}

/*
	Main functions
*/

func (p *Player) GiveClue(a *Action, g *Game) {
	// Keep track that someone clued (i.e. doing 1 clue costs 1 "Clue Token")
	g.Clues--

	// Apply the positive and negative clues to the cards in the hand
	for _, c := range p.Hand {
		positive := false
		if variantIsCardTouched(g, a.Clue, c) {
			positive = true
		}
		c.Clues = append(c.Clues, &CardClue{
			Type:     a.Clue.Type,
			Value:    a.Clue.Value,
			Positive: positive,
		})
		c.JustTouched = positive

		if a.Clue.Type == clueTypeRank {
			clueRank := a.Clue.Value
			for i := len(c.PossibleRanks) - 1; i >= 0; i-- {
				rank := c.PossibleRanks[i]
				if !(rank == clueRank == positive) {
					c.PossibleRanks = append(c.PossibleRanks[:i], c.PossibleRanks[i+1:]...)

					for _, suit := range variants[g.Variant].Suits {
						c.RemovePossibility(suit, rank, true)
					}
				}
			}
		} else if a.Clue.Type == clueTypeColor {
			clueSuit := variants[g.Variant].Suits[a.Clue.Value]
			for i := len(c.PossibleRanks) - 1; i >= 0; i-- {
				suit := c.PossibleSuits[i]
				if !(suit == clueSuit == positive) {
					c.PossibleSuits = append(c.PossibleSuits[:i], c.PossibleSuits[i+1:]...)

					for _, rank := range variants[g.Variant].Ranks {
						c.RemovePossibility(suit, rank, true)
					}
				}
			}
		}

		if len(c.PossibleSuits) == 1 && len(c.PossibleRanks) == 1 {
			c.Revealed = true
		}
	}

	p2 := g.Players[a.Target]
	str := "Turn " + strconv.Itoa(g.Turn+1) + " - " + p.Name + " tells " + p2.Name + " about "
	if a.Clue.Type == clueTypeRank {
		str += "rank " + strconv.Itoa(a.Clue.Value)
	} else if a.Clue.Type == clueTypeColor {
		color := variants[g.Variant].ClueColors[a.Clue.Value]
		str += "color " + color
	}
	log.Info(str)
	log.Info("There are now " + strconv.Itoa(g.Clues) + " clues left.")
}

func (p *Player) RemoveCard(target int, g *Game) *Card {
	// Get the target card
	i := p.GetCardIndex(target)
	c := p.Hand[i]

	// Remove it from the hand
	p.Hand = append(p.Hand[:i], p.Hand[i+1:]...)

	return c
}

func (p *Player) PlayCard(g *Game, c *Card) {
	// Find out if this successfully plays
	if c.IsPlayable(g) {
		c.Played = true
		g.Score++
		g.Stacks[c.Suit.Index] = c.Rank

		// Give the team a clue if the final card of the suit was played
		if c.Rank == 5 {
			g.Clues++

			// The extra clue is wasted if the team is at the maximum amount of clues already
			clueLimit := maxClues
			if g.Clues > clueLimit {
				g.Clues = clueLimit
			}
		}

		log.Info("Turn " + strconv.Itoa(g.Turn+1) + " - " + p.Name + " plays " + c.Name() + ".")
	} else {
		c.Failed = true
		p.DiscardCard(g, c)
		g.Strikes++
		log.Info("Turn " + strconv.Itoa(g.Turn+1) + " - " + p.Name + " fails to plays " + c.Name() + ". The team is now at " + strconv.Itoa(g.Strikes) + " strikes.")
	}
}

func (p *Player) DiscardCard(g *Game, c *Card) {
	c.Discarded = true
	g.DiscardPile = append(g.DiscardPile, c)
	if !c.Failed {
		log.Info("Turn " + strconv.Itoa(g.Turn+1) + " - " + p.Name + " discards " + c.Name() + ".")
	}
}

func (p *Player) DrawCard(g *Game) {
	// Don't draw any more cards if the deck is empty
	if g.DeckIndex >= len(g.Deck) {
		return
	}

	// Put it in the player's hand
	c := g.Deck[g.DeckIndex]
	g.DeckIndex++
	p.Hand = append(p.Hand, c)

	// Check to see if that was the last card drawn
	if g.DeckIndex >= len(g.Deck) {
		// Mark the turn upon which the game will end
		g.EndTurn = g.Turn + len(g.Players) + 1
	}

	log.Info("Turn " + strconv.Itoa(g.Turn+1) + " - " + p.Name + " draws a " + c.Name())
}

/*
	Subroutines
*/

func (p *Player) InHand(order int) bool {
	for _, c := range p.Hand {
		if c.Order == order {
			return true
		}
	}

	return false
}

func (p *Player) GetCardIndex(order int) int {
	for i, c := range p.Hand {
		if c.Order == order {
			return i
		}
	}

	return -1
}

func (p *Player) GetSlot(slot int) *Card {
	// Slot 1 is the newest (left-most) card, which is at index 4 (in a 3 player game)
	i := len(p.Hand) - slot
	if i < 0 || i > len(p.Hand)-1 {
		return nil
	}
	return p.Hand[i]
}

func (p *Player) GetCardsTouchedByClue(g *Game, clue *Clue) []*Card {
	touchedCards := make([]*Card, 0)
	for _, c := range p.Hand {
		if variantIsCardTouched(g, clue, c) {
			touchedCards = append(touchedCards, c)
		}
	}
	return touchedCards
}

func (p *Player) GetFreshCardsTouchedByClue(g *Game, clue *Clue) []*Card {
	freshCards := make([]*Card, 0)
	for _, c := range p.Hand {
		if c.IsClued() {
			continue
		}
		if variantIsCardTouched(g, clue, c) {
			freshCards = append(freshCards, c)
		}
	}
	return freshCards
}
