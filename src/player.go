package main

import (
	"sort"
	"strconv"
)

type Player struct {
	Index    int
	Name     string
	Hand     []*Card
	Notes    []string
	Strategy *Strategy
}

/*
	Main functions
*/

func (p *Player) GiveClue(a *Action, g *Game) {
	// Keep track that someone clued (performing a clue costs one "Clue Token")
	g.Clues--

	// Apply the positive and negative clues to the cards in the hand
	p2 := g.Players[a.Target]
	cardsTouched := 0
	slots := make([]int, 0)
	for _, c := range p2.Hand {
		positive := false
		if variantIsCardTouched(g, a.Clue, c) {
			c.Touched = true
			cardsTouched++
			slots = append(slots, c.Slot)
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
			for i := len(c.PossibleSuits) - 1; i >= 0; i-- {
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

	// Log the clue
	text := "Turn " + strconv.Itoa(g.Turn+1) + " - "
	text += p.Name + " tells " + p2.Name + " about "
	words := []string{
		"zero",
		"one",
		"two",
		"three",
		"four",
		"five",
	}
	text += words[cardsTouched] + " "
	if a.Clue.Type == clueTypeRank {
		text += strconv.Itoa(a.Clue.Value)
	} else if a.Clue.Type == clueTypeColor {
		text += variants[g.Variant].ClueColors[a.Clue.Value]
	}
	if cardsTouched != 1 {
		text += "s"
	}
	if cardsTouched > 0 {
		text += " in slot"
		if cardsTouched != 1 {
			text += "s"
		}
		text += " "
		sort.Ints(slots)
		for _, slot := range slots {
			text += strconv.Itoa(slot) + ", "
		}
		text = trimSuffix(text, ", ")
	}
	text += ". (There are now " + strconv.Itoa(g.Clues) + " clues left.)"
	log.Info(text)
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
		log.Info("Turn " + strconv.Itoa(g.Turn+1) + " - " + p.Name + " fails to plays " + c.Name() + ". (The team is now at " + strconv.Itoa(g.Strikes) + " strikes.)")
	}
	c.Holder = -1
	c.Slot = -1
}

func (p *Player) DiscardCard(g *Game, c *Card) {
	c.Discarded = true
	g.DiscardPile = append(g.DiscardPile, c)
	if !c.Failed {
		log.Info("Turn " + strconv.Itoa(g.Turn+1) + " - " + p.Name + " discards " + c.Name() + ".")
	}
	c.Holder = -1
	c.Slot = -1
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
	c.Holder = p.Index

	// Update the slot numbers of all of the cards in the hand
	for i, c2 := range p.Hand {
		c2.Slot = len(p.Hand) - i
	}

	// Check to see if that was the last card drawn
	if g.DeckIndex >= len(g.Deck) {
		// Mark the turn upon which the game will end
		g.EndTurn = g.Turn + len(g.Players) + 1
	}

	log.Info("Turn " + strconv.Itoa(g.Turn+1) + " - " + p.Name + " draws a " + c.Name() + ".")
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
		if c.Touched {
			continue
		}
		if variantIsCardTouched(g, clue, c) {
			freshCards = append(freshCards, c)
		}
	}
	return freshCards
}

func (p *Player) GetNextPlayerIndex(g *Game) int {
	return (p.Index + 1) % len(g.Players)
}

func (p *Player) GetPreviousPlayerIndex(g *Game) int {
	// In Golang, "%" will give the remainder and not the modulus,
	// so we need to ensure that the result is not negative or we will get a "index out of range" error
	return (p.Index - 1 + len(g.Players)) % len(g.Players)
}
