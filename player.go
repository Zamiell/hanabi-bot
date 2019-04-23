package main

import "strconv"

type Player struct {
	Index    int
	Name     string
	Hand     []*Card
	Strategy string
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

	log.Debug("Turn " + strconv.Itoa(g.Turn) + " - " + p.Name + " draws a " + c.Name(g))
}
