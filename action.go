package main

import (
	"strconv"
)

type Action struct {
	Type int `json:"type"` // 0 is play, 1 is discard, 2 is color clue, 3 is rank clue
	// If a clue, matches the index of the player
	// If a play/discard, matches the order of the card
	Target int `json:"target"`
	Value  int `json:"value"`
}

func actionClue(g *Game, p *Player, a *Action) {
	// Validate that the target of the clue is sane.
	if a.Target < 0 || a.Target > len(g.Players)-1 {
		logger.Fatal("The strategy of \"" + p.Strategy.Name + "\" " +
			"tried to clue an invalid clue target of \"" + strconv.Itoa(a.Target) + "\".")
		return
	}

	// Validate that the player is not giving a clue to themselves.
	if g.ActivePlayer == a.Target {
		logger.Fatal("The strategy of \"" + p.Strategy.Name + "\" " +
			"tried to give a clue to themself.")
		return
	}

	// Validate that there are clues available to use.
	if g.ClueTokens == 0 {
		logger.Fatal("The strategy of \"" + p.Strategy.Name + "\" " +
			"tried to give a clue while the team was at 0 clues.")
		return
	}

	// Validate that rank clues are valid.
	if a.Type == ActionTypeRankClue {
		valid := false
		for _, rank := range variants[g.Variant].ClueRanks {
			if rank == a.Value {
				valid = true
				break
			}
		}
		if !valid {
			logger.Fatal("The strategy of \"" + p.Strategy.Name + "\" tried to give a clue " +
				"with an invalid rank of \"" + strconv.Itoa(a.Value) + "\".")
			return
		}
	}

	// Validate that the color clues are valid.
	if a.Type == ActionTypeColorClue &&
		(a.Value < 0 || a.Value > len(variants[g.Variant].ClueColors)-1) {

		logger.Fatal("The strategy of \"" + p.Strategy.Name + "\" " +
			"tried to give a clue with an invalid color of \"" + strconv.Itoa(a.Value) + "\".")
		return
	}

	// Validate that the clue touches at least one card in the hand.
	touchedAtLeastOneCard := false
	p2 := g.Players[a.Target]
	for _, c := range p2.Hand {
		if variantIsCardTouched(g, NewClue(a), c) {
			touchedAtLeastOneCard = true
			break
		}
	}
	if !touchedAtLeastOneCard {
		logger.Fatal("The strategy of \"" + p.Strategy.Name + "\" tried to give a clue " +
			"that touches no cards in the hand.")
		return
	}

	p.GiveClue(a, g)
}

func actionPlay(g *Game, p *Player, a *Action) {
	// Validate that the card is in their hand.
	if !p.InHand(a.Target) {
		logger.Fatal("The strategy of \"" + p.Strategy.Name + "\" tried to play a card that was " +
			"not in their hand.")
		return
	}

	c := p.RemoveCard(a.Target)
	p.PlayCard(g, c)
	p.DrawCard(g)
}

func actionDiscard(g *Game, p *Player, a *Action) {
	// Validate that the card is in their hand.
	if !p.InHand(a.Target) {
		logger.Fatal("The strategy of \"" + p.Strategy.Name + "\" tried to discard a card that was " +
			"not in their hand.")
		return
	}

	// Validate that the team is not at the maximum amount of clues.
	// (The client should enforce this, but do a check just in case.)
	if g.ClueTokens == MaxClueNum {
		logger.Fatal("The strategy of \"" + p.Strategy.Name + "\" tried to discard while the team " +
			"was at the maximum amount of clues.")
		return
	}

	g.ClueTokens++
	c := p.RemoveCard(a.Target)
	p.DiscardCard(g, c)
	p.DrawCard(g)
}
