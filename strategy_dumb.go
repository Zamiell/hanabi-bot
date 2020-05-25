/*
	The "Dumb" strategy will alternate between blind-playing slot 1 and giving a clue.
	It is just intended to be a reference strategy.
*/

package main

// All strategies require a constructor that returns a reference.
func NewDumb() *Strategy {
	return &Strategy{
		Name: "Dumb",

		// Each strategy is composed of callback functions for various points in the game.
		Start: DumbStart, // Before the first move occurs
		// After a player announces what action they will perform but before it occurs
		ActionAnnounced: DumbActionAnnounced,
		ActionHappened:  DumbActionHappened, // After a player performs their action
		// When it gets to our turn; returns the action that we will perform
		GetAction: DumbGetAction,

		// A strategy's "data" can include both variables for game state and extra functions.
		Data: &Dumb{},
	}
}

type Dumb struct {
	Us        int  // Our player index
	BlindPlay bool // True if the player will blind-play on the next turn
}

// DumbStart is called before the first move occurs.
func DumbStart(s *Strategy, g *Game, us int) {
	d := s.Data.(*Dumb)

	// Store which player we are.
	d.Us = us

	// We don't want to blind-play on the first turn.
	d.BlindPlay = true
}

// DumbActionAnnounced is called after a player announces what action they will perform but before
// it occurs.
func DumbActionAnnounced(s *Strategy, g *Game, a *Action) {
	// We don't need to do anything at this step,
	// but a more complicated strategy would need to see "see" the new cards "touched" by a clue,
	// the slot number of a card that is going to be played, and so forth.
}

// DumbActionHappened is called after a player takes an action.
func DumbActionHappened(s *Strategy, g *Game, a *Action) {
	// We don't need to do anything at this step,
	// but a more complicated strategy would need to update internal variables relating to new
	// cards that are drawn, and so forth.
}

// DumbGetAction is called when it gets to our turn. It returns the action that we will perform.
func DumbGetAction(s *Strategy, g *Game) *Action {
	d := s.Data.(*Dumb)

	// Alternate between two brainless strategies.
	d.InvertDumbness()

	// Brainless strategy #1 - blind-play our slot 1 card.
	if d.BlindPlay && g.ClueTokens > 0 {
		// Get our slot 1 card
		p := g.Players[d.Us]
		firstCard := p.GetSlot(1)

		return &Action{
			Type:   ActionTypePlay,
			Target: firstCard.Order,
		}
	}

	// Brainless strategy #2 - give a rank clue to the next player's slot 1 card.
	// Get the next player.
	target := g.ActivePlayer + 1
	if target >= len(g.Players) {
		target = 0
	}

	// Get their slot 1 card.
	p := g.Players[target]
	firstCard := p.GetSlot(1)

	return &Action{
		Type:   ActionTypeRankClue,
		Target: target,
		Value:  firstCard.Rank,
	}
}

// We can attach new functions to the Dumb struct.
func (d *Dumb) InvertDumbness() {
	d.BlindPlay = !d.BlindPlay
}
