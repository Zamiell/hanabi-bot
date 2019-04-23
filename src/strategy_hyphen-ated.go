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
	OurIndex  int
	Deck      []*HyphenCard
	EarlyGame bool
}
type HyphenCard struct {
	KnownPlayable bool
	KnownTrash    bool
}

// HyphenatedStart is called before the first move occurs
func HyphenatedStart(s *Strategy, g *Game, us int) {
	d := s.Data.(*Hyphenated)

	// Store which player we are
	d.OurIndex = us

	// Make a copy of the deck that will store additional information about each card
	d.Deck = make([]*HyphenCard, 0)
	for i := 0; i < len(g.Deck); i++ {
		d.Deck = append(d.Deck, &HyphenCard{})
	}
}

// HyphenatedGetAction is called when it gets to our turn
// It returns the action that we will perform
func HyphenatedGetAction(s *Strategy, g *Game) *Action {
	d := s.Data.(*Hyphenated)
	var a *Action

	if g.Clues > 0 {
		// Check for the next guy's chop

		// Check for playable cards

		// Clue playable cards
		a = d.LookFor1s(g)
		if a != nil {
			return a
		}
	}

	a = d.CheckPlayable(g)
	if a != nil {
		return a
	}

	return nil
}

// HyphenatedActionHappened is called when an action happens
func HyphenatedActionHappened(s *Strategy, g *Game, a *Action) {
	// d := s.Data.(*Hyphenated)

	if a.Type == actionTypeClue {

	}
}

func (d *Hyphenated) CheckPlayable(g *Game) *Action {
	// Look through our whole hand and make a list of all the playable cards
	// TODO
	return nil
}

func (d *Hyphenated) LookFor1s(g *Game) *Action {
	// TODO
	return nil
}
