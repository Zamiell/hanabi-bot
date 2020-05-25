/*
	The "Hyphen-ated" strategy is intended to roughly match the strategies for the Hyphen-ated
	group.
	https://github.com/Zamiell/hanabi-conventions/blob/master/Reference.md
*/

package main

import (
	"strconv"
)

func NewHyphenated() *Strategy {
	return &Strategy{
		Name:            "Hyphen-ated",
		Start:           HyphenatedStart,
		GetAction:       HyphenatedGetAction,
		ActionAnnounced: HyphenatedActionAnnounced,
		ActionHappened:  HyphenatedActionHappened,
		Data:            &Hyphenated{},
	}
}

type Hyphenated struct {
	Us        int // Our index
	Players   []*HyphenPlayer
	Cards     []*HyphenCard
	EarlyGame bool
}

// HyphenatedStart is called before the first move occurs.
func HyphenatedStart(s *Strategy, g *Game, us int) {
	d := s.Data.(*Hyphenated)

	// Store which player we are.
	d.Us = us

	// Initialize the objects that will store additional information about each player.
	d.Players = make([]*HyphenPlayer, 0)
	for i := 0; i < len(g.Players); i++ {
		d.Players = append(d.Players, &HyphenPlayer{
			Index: i,
		})
	}

	// Initialize the objects that will store additional information about each card.
	d.Cards = make([]*HyphenCard, 0)
	for i := 0; i < len(g.Deck); i++ {
		d.Cards = append(d.Cards, &HyphenCard{
			Order: i,
			// KnownPossibleCards: make([]string),
		})
	}
}

// HyphenatedActionAnnounced is called before a player clues, plays, or discards.
func HyphenatedActionAnnounced(s *Strategy, g *Game, a *Action) {
	d := s.Data.(*Hyphenated)

	if a.Type == ActionTypeColorClue || a.Type == ActionTypeRankClue {
		HyphenatedActionAnnouncedClue(d, g, a)
	}
}

// HyphenatedActionHappened is called after a player clues, plays, or discards.
func HyphenatedActionHappened(s *Strategy, g *Game, a *Action) {
	d := s.Data.(*Hyphenated)

	if a.Type == ActionTypeColorClue || a.Type == ActionTypeRankClue {
		HyphenatedActionHappenedClue(d, g, a)
	}
}

// HyphenatedGetAction is called when it gets to our turn.
// It returns the action that we will perform.
func HyphenatedGetAction(s *Strategy, g *Game) *Action {
	d := s.Data.(*Hyphenated)
	var a *Action

	n := 0
	if g.ClueTokens > 0 {
		// Check to see if the next player has a safe discard.
		a = d.CheckNextPlayerSave(g)
		n++
		if a != nil {
			logger.Info("Using logic " + strconv.Itoa(n) + ": CheckNextPlayerSave")
			return a
		}

		// If we have playable cards, play them.
		a = d.CheckPlayableCards(g)
		n++
		if a != nil {
			logger.Info("Using logic " + strconv.Itoa(n) + ": CheckPlayableCards")
			return a
		}

		// Clue playable cards.
		a = d.CheckPlayClues(g)
		n++
		if a != nil {
			logger.Info("Using logic " + strconv.Itoa(n) + ": CheckPlayClues")
			return a
		}
	}

	// If we have playable cards, play them.
	a = d.CheckPlayableCards(g)
	n++
	if a != nil {
		logger.Info("Using logic " + strconv.Itoa(n) + ": CheckPlayableCards")
		return a
	}

	if g.ClueTokens != MaxClueNum {
		a = d.Discard(g)
		n++
		if a != nil {
			logger.Info("Using logic " + strconv.Itoa(n) + ": Discard")
			return a
		}
	}

	if g.ClueTokens > 0 {
		// Give a 5 Stall clue.
		a = d.Check5Stall(g)
		n++
		if a != nil {
			logger.Info("Using logic " + strconv.Itoa(n) + ": Check5Stall")
			return a
		}

		// Give a "hard burn" (e.g. a stall clue).
		a = d.Check5Burn(g)
		n++
		if a != nil {
			logger.Info("Using logic " + strconv.Itoa(n) + ": Check5Burn")
			return a
		}
	}

	return nil
}
