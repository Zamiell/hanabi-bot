package main

import (
	"strconv"
)

type Player struct {
	Index    int
	Name     string
	Hand     []*Card
	Strategy *Strategy
}

func (p *Player) GiveClue(a *Action, g *Game) {
	// Apply the positive and negative clues
	for _, c := range p.Hand {
		positive := false
		if variantIsCardTouched(g.Variant, a.Clue, c) {
			positive = true
		}
		c.Clues = append(c.Clues, &CardClue{
			Type:     a.Clue.Type,
			Value:    a.Clue.Value,
			Positive: positive,
		})
	}

	// Keep track that someone clued (i.e. doing 1 clue costs 1 "Clue Token")
	g.Clues--

	p2 := g.Players[a.Target]
	str := "Turn " + strconv.Itoa(g.Turn+1) + " - " + p.Name + " tells " + p2.Name + " about "
	if a.Clue.Type == clueTypeRank {
		str += "rank " + strconv.Itoa(a.Clue.Value)
	} else if a.Clue.Type == clueTypeColor {
		color := variants[g.Variant].ClueColors[a.Clue.Value]
		str += "color " + color
	}
	log.Info(str)
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
	if c.Rank == g.Stacks[c.Suit]+1 {
		g.Score++
		g.Stacks[c.Suit] = c.Rank

		// Give the team a clue if the final card of the suit was played
		if c.Rank == 5 {
			g.Clues++

			// The extra clue is wasted if the team is at the maximum amount of clues already
			clueLimit := maxClues
			if g.Clues > clueLimit {
				g.Clues = clueLimit
			}
		}

		log.Info("Turn " + strconv.Itoa(g.Turn+1) + " - " + p.Name + " plays " + c.Name(g) + ".")
	} else {
		g.Strikes++
		log.Info("Turn " + strconv.Itoa(g.Turn+1) + " - " + p.Name + " fails to plays " + c.Name(g) + ". The team is now at " + strconv.Itoa(g.Strikes) + " strikes.")
	}
}

func (p *Player) DiscardCard(g *Game, c *Card) {
	g.DiscardPile = append(g.DiscardPile, c)
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

	log.Debug("Turn " + strconv.Itoa(g.Turn+1) + " - " + p.Name + " draws a " + c.Name(g))
}

/*
	Subroutines
*/

func (p *Player) GetCardIndex(order int) int {
	for i, c := range p.Hand {
		if c.Order == order {
			return i
		}
	}

	return -1
}

func (p *Player) InHand(order int) bool {
	for _, c := range p.Hand {
		if c.Order == order {
			return true
		}
	}

	return false
}
