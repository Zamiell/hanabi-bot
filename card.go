package main

import (
	"strconv"
)

type Card struct {
	Suit          *Suit
	Rank          int
	Order         int
	Holder        int
	Slot          int
	Clues         []*CardClue // This is a list that includes both positive and negative clues
	Touched       bool        // True if it has one or more positive clues on it
	JustTouched   bool        // Touched by the last clue that was given
	PossibleSuits []*Suit
	PossibleRanks []int
	PossibleCards map[string]int // Maps card identities to count
	Revealed      bool
	Played        bool
	Discarded     bool
	Failed        bool // If the card failed to play
}

func (c *Card) Name() string {
	name := c.Suit.Name + " " + strconv.Itoa(c.Rank)
	return name
}

func (c *Card) IsPlayable(g *Game) bool {
	return c.Rank == g.Stacks[c.Suit.Index]+1
}

func (c *Card) IsCritical(g *Game) bool {
	if c.Played || c.Discarded {
		return false
	}

	total, discarded := g.GetSpecificCardNum(c.Suit, c.Rank)
	return total == discarded+1
}

func (c *Card) CouldBeSuit(s *Suit) bool {
	for _, suit := range c.PossibleSuits {
		if s == suit {
			return true
		}
	}
	return false
}

func (c *Card) CouldBeRank(r int) bool {
	for _, rank := range c.PossibleRanks {
		if r == rank {
			return true
		}
	}
	return false
}

// NeedsToBePlayed returns true if the card is not yet played
// and is still needed to be played in order to get the maximum score
func (c *Card) NeedsToBePlayed(g *Game) bool {
	// First, check to see if a copy of this card has already been played
	for _, c2 := range g.Deck {
		if c2.Suit == c.Suit &&
			c2.Rank == c.Rank &&
			c2.Played {

			return false
		}
	}

	// Second, check to see if it is still possible to play this card
	// (the preceding cards in the suit might have already been discarded)
	for i := 1; i < c.Rank; i++ {
		total, discarded := g.GetSpecificCardNum(c.Suit, i)
		if total == discarded {
			// The suit is "dead", so this card does not need to be played anymore
			return false
		}
	}

	// By default, all cards not yet played will need to be played
	return true
}

func (c *Card) RemovePossibility(suit *Suit, rank int, removeAll bool) {
	// Every card has a possibility map that maps card identities to count
	mapIndex := suit.Name + strconv.Itoa(rank)
	cardsLeft := c.PossibleCards[mapIndex]
	if cardsLeft > 0 {
		// Remove one or all possibilities for this card,
		// (depending on whether the card was clued
		// or if we saw someone draw a copy of this card)
		cardsLeft := cardsLeft - 1
		if removeAll {
			cardsLeft = 0
		}
		c.PossibleCards[mapIndex] = cardsLeft
	}
}
