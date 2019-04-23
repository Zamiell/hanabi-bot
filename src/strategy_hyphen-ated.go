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
	BlindPlay bool // True if the player will blind-play on the next turn
}

// HyphenatedStart is called before the first move occurs
func HyphenatedStart(s *Strategy) {
	// d := s.Data.(*Dumb)
}

// HyphenatedGetAction is called when it gets to our turn
// It returns the action that we will perform
func HyphenatedGetAction(s *Strategy, g *Game) *Action {
	// d := s.Data.(*Dumb)
	return nil
}

// HyphenatedActionHappened is called when an action happens
func HyphenatedActionHappened(s *Strategy, g *Game) {
}

// We can attach new functions to the Dumb struct
func (d *Hyphenated) DoSomething() {
}
