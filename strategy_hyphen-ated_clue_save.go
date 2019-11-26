package main

func (d *Hyphenated) CheckNextPlayerSave(g *Game) *Action {
	p := g.Players[d.Us]
	npi := p.GetNextPlayerIndex(g)
	np := g.Players[npi]
	hnp := d.Players[npi]
	c := hnp.GetChop(g, d)
	hc := d.Cards[c.Order]

	// Look to see if the next player's chop card needs to be given a Play Clue or a Save Clue
	if !hc.Known &&
		!hc.Playable &&
		!hc.Trash &&
		!hc.ChopMoved &&
		!c.IsCritical(g) &&
		!hc.IsAlreadyTouchedOrChopMoved(g, d) &&
		!hc.IsPlayableUnique(g, d) &&
		!hc.Needs2Save(g, d) {

		return nil
	}
	if c.IsCritical(g) {
		log.Debug(c.Name() + " on " + np.Name + "'s chop is critical.")
	} else if hc.IsPlayableUnique(g, d) {
		log.Debug(c.Name() + " on " + np.Name + "'s chop is playable and unique.")
	} else if hc.Needs2Save(g, d) {
		log.Debug(c.Name() + " on " + np.Name + "'s chop needs a 2 Save.")
	}

	// We need to make sure that this card is not discarded
	// First, check to see if we can give the card a Play Clue instead of a Save Clue
	// TODO: Instead of just looking at this specific card,
	// look through this player's entire hand for other playable cards,
	// since we should prefer to give Play Clues over Save Clues
	a := hc.CheckPlayClue(g, d)
	if a != nil {
		return a
	}

	// We cannot give a Play Clue to the important card, so we need to give a Save Clue
	if !c.IsCritical(g) || c.Rank == 5 {
		// We must use rank to save this card (e.g. 2 Save or 5 Save)
		return &Action{
			Type: actionTypeClue,
			Clue: &Clue{
				Type:  clueTypeRank,
				Value: c.Rank,
			},
			Target: npi,
		}
	}

	// We can use rank or color to save this card
	return hc.CheckRankOrColor(g, d, hyphenClueTypeSave)
}
