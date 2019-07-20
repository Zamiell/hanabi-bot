package main

// We store extra information about each card
type HyphenCard struct {
	Order         int
	PossibleCards []*PossibleCard
	Known         bool
	Playable      bool
	Trash         bool
	ChopMoved     bool
}

type PossibleCard struct {
	Suit *Suit
	Rank int
}

func (hc *HyphenCard) IsAlreadyTouchedOrChopMoved(g *Game, d *Hyphenated) bool {
	c := g.Deck[hc.Order]

	for i, p := range g.Players {
		for _, c2 := range p.Hand {
			// Skip over the card we are looking for duplicates of
			if c.Order == c2.Order {
				continue
			}

			// Skip over cards that are not copies
			if c.Suit != c2.Suit || c.Rank != c2.Rank {
				continue
			}

			hc2 := d.Cards[c2.Order]
			if i == d.Us {
				if hc2.Known {
					return true
				}
			} else {
				if c2.Touched || hc2.Known || hc2.ChopMoved {
					return true
				}
			}
		}
	}

	return false
}

// Is this card both playable and the only visible copy?
func (hc *HyphenCard) IsPlayableUnique(g *Game, d *Hyphenated) bool {
	c := g.Deck[hc.Order]

	if c.Rank <= g.Stacks[c.Suit.Index] {
		return false
	}

	if !c.IsPlayable(g) && !hc.IsDelayedPlayable(g, d) {
		return false
	}

	return !hc.IsCopyInOtherHand(g, d)
}

// Check to see if this card needs to be given a 2 Save clue
func (hc *HyphenCard) Needs2Save(g *Game, d *Hyphenated) bool {
	c := g.Deck[hc.Order]

	if c.Rank != 2 {
		return false
	}

	// If the card is critical, this will be detected elsewhere
	// If the card is playable, this will be detected elsewhere
	// If the card is already touched or chop moved, this will be detected elsewhere
	return !hc.IsCopyInOtherHand(g, d)
}

// Is this card playable through a clued card in someone else's hand?
func (hc *HyphenCard) IsDelayedPlayable(g *Game, d *Hyphenated) bool {
	c := g.Deck[hc.Order]

	if c.Rank <= g.Stacks[c.Suit.Index]+1 {
		return false
	}

	// Check to see if all the cards leading up to this card are clued in someone's hand
	for inBetweenRank := g.Stacks[c.Suit.Index] + 1; inBetweenRank < c.Rank; inBetweenRank++ {
		cardTouched := false

		for i, p := range g.Players {
			for _, c2 := range p.Hand {
				// Skip over cards that are not copies
				if c.Suit != c2.Suit || inBetweenRank != c2.Rank {
					continue
				}

				hc2 := d.Cards[c2.Order]
				if i == d.Us {
					// Check to see if the card is touched in our hand
					// We want to check for "Known" instead of "Touched"
					// since we want this to include Finessed cards
					if hc2.Known {
						cardTouched = true
						break
					}
				} else {
					// Check to see if the card is touched in anyone else's hand
					if c2.Touched || hc2.Known {
						cardTouched = true
						break
					}
				}
			}
		}

		if !cardTouched {
			return false
		}
	}

	return true
}

func (hc *HyphenCard) IsCopyInOtherHand(g *Game, d *Hyphenated) bool {
	c := g.Deck[hc.Order]

	for i, p := range g.Players {
		for _, c2 := range p.Hand {
			// Skip over the card we are looking for duplicates of
			if c.Order == c2.Order {
				continue
			}

			// Skip over cards that are not copies
			if c.Suit != c2.Suit || c.Rank != c2.Rank {
				continue
			}

			hc2 := d.Cards[c2.Order]
			if i == d.Us {
				if hc2.Known {
					return true
				}
			} else {
				return true
			}
		}
	}

	return false
}

// Check to see if this card can be given a Play Clue
func (hc *HyphenCard) CheckPlayClue(g *Game, d *Hyphenated) *Action {
	c := g.Deck[hc.Order]

	if !c.IsPlayable(g) && !hc.IsDelayedPlayable(g, d) {
		return nil
	}

	return hc.CheckRankOrColor(g, d, hyphenClueTypePlay)
}

// Check to see if a rank clue or a color clue would be better
// (this function works for both Play Clues and Save Clues)
func (hc *HyphenCard) CheckRankOrColor(g *Game, d *Hyphenated, hyphenClueType int) *Action {
	c := g.Deck[hc.Order]

	rankClue := d.CheckViableClue(g, c.Holder, clueTypeRank, c.Rank, hyphenClueType)
	colorClue := d.CheckViableClue(g, c.Holder, clueTypeColor, c.Suit.ColorValue(g), hyphenClueType)

	if rankClue == nil && colorClue == nil {
		// Both a rank Play Clue and a color Play Clue would not be viable for some reason
		// (e.g. it would violate Good Touch Principle)
		return nil
	}

	var clueType int
	if colorClue == nil {
		// A color clue would duplicate something but a rank clue would be clean
		clueType = clueTypeRank
	} else if rankClue == nil {
		// A rank clue would duplicate something but a color clue would be clean
		clueType = clueTypeColor
	} else if rankClue.CardsClued > colorClue.CardsClued {
		// Both clues are clean but a rank clue would touch more useful cards
		clueType = clueTypeRank
	} else {
		// Prefer a color clue over a rank clue if the amount of cards touched is equal or greater
		clueType = clueTypeColor
	}

	var clueValue int
	if clueType == clueTypeRank {
		clueValue = c.Rank
	} else if clueType == clueTypeColor {
		clueValue = c.Suit.ColorValue(g)
	}

	return &Action{
		Type: actionTypeClue,
		Clue: &Clue{
			Type:  clueType,
			Value: clueValue,
		},
		Target: c.Holder,
	}
}
