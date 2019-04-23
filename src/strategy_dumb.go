/*
	The "Dumb" strategy will alternate between blind-playing slot 1 and giving a clue
	It is just intended to be a reference strategy
*/

package main

func NewDumb() *Strategy {
	return &Strategy{
		Name:           "Dumb",
		Start:          DumbStart,
		GetAction:      DumbGetAction,
		ActionHappened: DumbActionHappened,
		Data:           &Dumb{},
	}
}

type Dumb struct {
	Us        int  // Our player index
	BlindPlay bool // True if the player will blind-play on the next turn
}

// DumbStart is called before the first move occurs
func DumbStart(s *Strategy, g *Game, us int) {
	d := s.Data.(*Dumb)

	// Store which player we are
	d.Us = us

	// We don't want to blind-play on the first turn
	d.BlindPlay = true
}

// DumbActionHappened is called when an action happens
func DumbActionHappened(s *Strategy, g *Game, a *Action) {
}

// DumbGetAction is called when it gets to our turn
// It returns the action that we will perform
func DumbGetAction(s *Strategy, g *Game) *Action {
	d := s.Data.(*Dumb)

	// Alternate between two brainless strategies
	d.InvertDumbness()

	if d.BlindPlay && g.Clues > 0 {
		// Get our slot 1 card
		p := g.Players[d.Us]
		firstCard := p.GetSlot(1)

		return &Action{
			Type:   actionTypePlay,
			Target: firstCard.Order,
		}
	}

	// Get the next player
	target := g.ActivePlayer + 1
	if target >= len(g.Players) {
		target = 0
	}

	// Get their slot 1 card
	p := g.Players[target]
	firstCard := p.GetSlot(1)

	return &Action{
		Type: actionTypeClue,
		Clue: &Clue{
			Type:  clueTypeRank,
			Value: firstCard.Rank,
		},
		Target: target,
	}
}

// We can attach new functions to the Dumb struct
func (d *Dumb) InvertDumbness() {
	d.BlindPlay = !d.BlindPlay
}
