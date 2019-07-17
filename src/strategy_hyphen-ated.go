/*
	The "Hyphen-ated" strategy is intended to roughly match the strategies for the Hyphen-ated group
	https://github.com/Zamiell/hanabi-conventions/blob/master/Reference.md
*/

package main

import "strconv"

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

// We store extra information about each card
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

	// Initialize the objects that will store additional information about each player
	d.Players = make([]*HyphenPlayer, 0)
	for i := 0; i < len(g.Players); i++ {
		d.Players = append(d.Players, &HyphenPlayer{
			Index: i,
		})
	}

	// Initialize the objects that will store additional information about each card
	d.Cards = make([]*HyphenCard, 0)
	for i := 0; i < len(g.Deck); i++ {
		d.Cards = append(d.Cards, &HyphenCard{})
	}
}

// HyphenatedActionHappened is called when a player clues, plays, or discards
func HyphenatedActionHappened(s *Strategy, g *Game, a *Action) {
	if a.Type == actionTypeClue {
		HyphenatedActionHappenedClue(s, g, a)
	} else if a.Type == actionTypePlay {
	} else if a.Type == actionTypeDiscard {
	}
}

func HyphenatedActionHappenedClue(s *Strategy, g *Game, a *Action) {
	d := s.Data.(*Hyphenated)
	p := g.Players[a.Target]
	hp := d.Players[a.Target]
	interpretation := d.GetClueInterpretation(g, a)
	if interpretation == hyphenClueTypePlay {
		focusedCard := d.GetClueFocus(g, a.Target, a.Clue)
		d.Cards[focusedCard.Order].Playable = true
	}

	// Update all playable cards based on good touch principle
	// TODO 1: this must also be done after a card is played
	// TODO 2: don't mark a card as playable if it is known trash
	for _, c := range p.Hand {
		if c.IsClued() && !d.Cards[c.Order].Playable {
			isplayable := true
			for _, s := range c.PossibleSuits {
				for _, r := range c.PossibleRanks {
					if r > g.Stacks[s.Index]+1 {
						isplayable = false
						break
					}
				}
			}
			if isplayable {
				// log.Info(c.Name() + " is playable." + strconv.FormatBool(c.IsClued()))
				d.Cards[c.Order].Playable = true
			}
		}
	}

	hp.UpdateChop(g)
}

// HyphenatedGetAction is called when it gets to our turn
// It returns the action that we will perform
func HyphenatedGetAction(s *Strategy, g *Game) *Action {
	d := s.Data.(*Hyphenated)
	var a *Action

	n := 0
	if g.Clues > 0 {
		// Check for the next guy's chop
		a = d.CheckNextPlayerChop(g)
		n++
		if a != nil {
			log.Info("Using logic " + strconv.Itoa(n) + ": CheckNextPlayerChop")
			return a
		}

		// If we have playable cards, play them
		a = d.CheckPlayableCards(g)
		n++
		if a != nil {
			log.Info("Using logic " + strconv.Itoa(n) + ": CheckPlayableCards")
			return a
		}

		// Clue playable cards
		a = d.CheckPlayClues(g)
		n++
		if a != nil {
			log.Info("Using logic " + strconv.Itoa(n) + ": CheckPlayClues")
			return a
		}
	}

	// If we have playable cards, play them
	a = d.CheckPlayableCards(g)
	n++
	if a != nil {
		log.Info("Using logic " + strconv.Itoa(n) + ": CheckPlayableCards")
		return a
	}

	if g.Clues != 8 {
		a = d.Discard(g)
		n++
		if a != nil {
			log.Info("Using logic " + strconv.Itoa(n) + ": Discard")
			return a
		}
	}

	if g.Clues > 0 {
		// Give a 5 Stall clue
		a = d.Check5Stall(g)
		n++
		if a != nil {
			log.Info("Using logic " + strconv.Itoa(n) + ": Check5Stall")
			return a
		}

		// Give a "hard burn" (e.g. a stall clue)
		a = d.Check5Burn(g)
		n++
		if a != nil {
			log.Info("Using logic " + strconv.Itoa(n) + ": Check5Burn")
			return a
		}
	}

	return nil
}
