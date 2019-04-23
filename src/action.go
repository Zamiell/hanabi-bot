package main

import (
	"strconv"
)

func actionClue(g *Game, p *Player, a *Action) {
	// Validate that the target of the clue is sane
	if a.Target < 0 || a.Target > len(g.Players)-1 {
		log.Fatal("The strategy of \"" + p.Strategy.Name + "\" tried to clue an " +
			"invalid clue target of \"" + strconv.Itoa(a.Target) + "\".")
		return
	}

	// Validate that the player is not giving a clue to themselves
	if g.ActivePlayer == a.Target {
		log.Fatal("The strategy of \"" + p.Strategy.Name + "\" tried to give a clue " +
			"to themself.")
		return
	}

	// Validate that there are clues available to use
	if g.Clues == 0 {
		log.Fatal("The strategy of \"" + p.Strategy.Name + "\" tried to give a clue " +
			"while the team was at 0 clues.")
		return
	}

	// Validate that the clue type is sane
	if a.Clue.Type < clueTypeRank || a.Clue.Type > clueTypeColor {
		log.Fatal("The strategy of \"" + p.Strategy.Name + "\" tried to give a clue " +
			"with an invalid clue type of \"" + strconv.Itoa(a.Clue.Type) + "\".")
		return
	}

	// Validate that rank clues are valid
	if a.Clue.Type == clueTypeRank {
		valid := false
		for _, rank := range variants[g.Variant].ClueRanks {
			if rank == a.Clue.Value {
				valid = true
				break
			}
		}
		if !valid {
			log.Fatal("The strategy of \"" + p.Strategy.Name + "\" tried to give a clue " +
				"with an invalid rank of \"" + strconv.Itoa(a.Clue.Value) + "\".")
			return
		}
	}

	// Validate that the color clues are valid
	if a.Clue.Type == clueTypeColor &&
		(a.Clue.Value < 0 || a.Clue.Value > len(variants[g.Variant].ClueColors)-1) {

		log.Fatal("The strategy of \"" + p.Strategy.Name + "\" tried to give a clue " +
			"with an invalid color of \"" + strconv.Itoa(a.Clue.Value) + "\".")
		return
	}

	// Validate that the clue touches at least one card in the hand
	touchedAtLeastOneCard := false
	p2 := g.Players[a.Target]
	for _, c := range p2.Hand {
		if variantIsCardTouched(g, a.Clue, c) {
			touchedAtLeastOneCard = true
			break
		}
	}
	if !touchedAtLeastOneCard {
		log.Fatal("The strategy of \"" + p.Strategy.Name + "\" tried to give a clue " +
			"that touches no cards in the hand.")
		return
	}

	p.GiveClue(a, g)
}

func actionPlay(g *Game, p *Player, a *Action) {
	// Validate that the card is in their hand
	if !p.InHand(a.Target) {
		log.Fatal("The strategy of \"" + p.Strategy.Name + "\" tried to play a card that was " +
			"not in their hand.")
		return
	}

	c := p.RemoveCard(a.Target, g)
	p.PlayCard(g, c)
	p.DrawCard(g)
}

func actionDiscard(g *Game, p *Player, a *Action) {
	// Validate that the card is in their hand
	if !p.InHand(a.Target) {
		log.Fatal("The strategy of \"" + p.Strategy.Name + "\" tried to discard a card that was " +
			"not in their hand.")
		return
	}

	// Validate that the team is not at the maximum amount of clues
	// (the client should enforce this, but do a check just in case)
	if g.Clues == maxClues {
		log.Fatal("The strategy of \"" + p.Strategy.Name + "\" tried to discard while the team " +
			"was at the maximum amount of clues.")
		return
	}

	g.Clues++
	c := p.RemoveCard(a.Target, g)
	p.DiscardCard(g, c)
	p.DrawCard(g)
}
