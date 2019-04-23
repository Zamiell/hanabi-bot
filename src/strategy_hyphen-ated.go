/*
	The "Hyphen-ated" strategy is intended to roughly match the strategies for the Hyphen-ated group
	https://github.com/Zamiell/hanabi-conventions/blob/master/Reference.md
*/

package main

func NewHyphenated() *Strategy {
	return &Strategy{
		Name:           "Hyphen-ated",
		Start:          HyphenatedStart,
		GetAction:      HyphenatedGetAction,
		ActionHappened: HyphenatedActionHappened,
		Data:           &Hyphenated{},
	}
}

type Hyphenated struct {
	Us        int // Our index
	Players   []*HyphenPlayer
	Cards     []*HyphenCard
	EarlyGame bool
}
type HyphenCard struct {
	Playable bool
	Trash    bool
}

const (
	hyphenClueTypeSave = iota
	hyphenClueTypePlay
	// hyphenClueTypeFix
)

// HyphenatedStart is called before the first move occurs
func HyphenatedStart(s *Strategy, g *Game, us int) {
	d := s.Data.(*Hyphenated)

	// Store which player we are
	d.Us = us

	// We need to store additional information about each player
	d.Players = make([]*HyphenPlayer, 0)
	for i := 0; i < len(g.Players); i++ {
		d.Players = append(d.Players, &HyphenPlayer{
			Index: i,
		})
	}

	// We need to store additional information about each card
	d.Cards = make([]*HyphenCard, 0)
	for i := 0; i < len(g.Deck); i++ {
		d.Cards = append(d.Cards, &HyphenCard{})
	}
}

// HyphenatedActionHappened is called when an action happens
func HyphenatedActionHappened(s *Strategy, g *Game, a *Action) {
	d := s.Data.(*Hyphenated)

	if a.Type == actionTypeClue {
		p := g.Players[a.Target]
		touchedCards := make([]*Card, 0)
		for _, c := range p.Hand {
			if c.JustTouched {
				touchedCards = append(touchedCards, c)
			}
		}

		interpretation := d.GetClueInterpretation(g, a)
		if interpretation == hyphenClueTypePlay {
			focusedCard := d.GetClueFocus(g, a.Target, a.Clue)
			d.Cards[focusedCard.Order].Playable = true
		}

		d.UpdateChop(g, a)
	}
}

// HyphenatedGetAction is called when it gets to our turn
// It returns the action that we will perform
func HyphenatedGetAction(s *Strategy, g *Game) *Action {
	d := s.Data.(*Hyphenated)
	var a *Action

	if g.Clues > 0 {
		// Check for the next guy's chop
		// TODO

		// Check for playable cards
		a = d.CheckPlayable(g)
		if a != nil {
			return a
		}

		// Clue playable cards
		a = d.CheckPlayClues(g)
		if a != nil {
			return a
		}
	}

	a = d.CheckPlayable(g)
	if a != nil {
		return a
	}

	if g.Clues != 8 {
		a = d.Discard(g)
		if a != nil {
			return a
		}
	}

	// TODO give a stall clue

	return nil
}
