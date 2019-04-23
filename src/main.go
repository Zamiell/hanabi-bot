package main

import (
	"math/rand"
	"os"
	"strconv"

	"github.com/op/go-logging"
)

var (
	log *logging.Logger
)

func main() {
	// Initialize logging
	// http://godoc.org/github.com/op/go-logging#Formatter
	log = logging.MustGetLogger("hanabi-bot")
	loggingBackend := logging.NewLogBackend(os.Stdout, "", 0)
	logFormat := logging.MustStringFormatter( // https://golang.org/pkg/time/#Time.Format
		`%{time:Mon Jan 02 15:04:05 MST 2006} - %{level:.4s} - %{shortfile} - %{message}`,
	)
	loggingBackendFormatted := logging.NewBackendFormatter(loggingBackend, logFormat)
	logging.SetBackend(loggingBackendFormatted)

	// Welcome message
	log.Info("+----------------------+")
	log.Info("| Starting hanabi-bot. |")
	log.Info("+----------------------+")

	variantsInit()
	stratInit()

	// Initialize the game
	g := &Game{
		Variant:     "No Variant",
		Players:     make([]*Player, 0),
		Stacks:      make([]int, 0),
		DiscardPile: make([]*Card, 0),
		Clues:       maxClues,
		Actions:     make([]*Action, 0),
		EndTurn:     -1,
	}

	// Initialize the players
	names := []string{"Alice", "Bob", "Cathy"}
	for _, name := range names {
		strategy := "Hyphen-ated"
		p := &Player{
			Name:     name,
			Hand:     make([]*Card, 0),
			Strategy: strategies[strategy](),
		}
		g.Players = append(g.Players, p)
	}

	g.InitDeck()
	g.InitStacks()
	rand.Seed(int64(g.Seed)) // Seed the random number generator with the game seed
	g.Shuffle()
	g.InitStartingPlayer()
	g.DealStartingHands()

	for _, p := range g.Players {
		p.Strategy.Start(p.Strategy)
	}

	// Play the game until it ends
	for {
		// Query the strategy to see what kind of move that the player will do
		p := g.Players[g.ActivePlayer]
		a := p.Strategy.GetAction(p.Strategy, g)
		if a == nil {
			log.Fatal("The strategy of \"" + p.Strategy.Name + "\" returned a nil action.")
		}

		// Perform the move
		if a.Type == actionTypeClue {
			actionClue(g, p, a)
		} else if a.Type == actionTypePlay {
			actionPlay(g, p, a)
		} else if a.Type == actionTypeDiscard {
			actionDiscard(g, p, a)
		} else {
			log.Fatal("The strategy of \"" + p.Strategy.Name + "\" returned an illegal action type of " +
				"\"" + strconv.Itoa(a.Type) + "\".")
			return
		}
		g.Actions = append(g.Actions, a)

		// Increment the turn
		g.Turn++
		g.ActivePlayer = (g.ActivePlayer + 1) % len(g.Players)
		if g.CheckEnd() {
			log.Info("----------------------------------------")
			if g.EndCondition > endConditionNormal {
				log.Info("Players lose!")
			} else {
				log.Info("Players score " + strconv.Itoa(g.Score) + " points.")
			}
		}

		if g.EndCondition > endConditionInProgress {
			break
		}
	}
}

func actionClue(g *Game, p *Player, a *Action) {
	// Validate that the target of the clue is sane
	if a.Target < 0 || a.Target > len(g.Players)-1 {
		log.Fatal("The strategy of \"" + p.Strategy.Name + "\" tried to clue an " +
			"invalid clue target of \"" + strconv.Itoa(a.Target) + "\".")
		return
	}

	// Validate that the player is not giving a clue to themselves
	if g.ActivePlayer == a.Target {
		log.Fatal("The strategy of \"" + p.Strategy.Name + "\" tried to give a clue " +
			"to themself.")
		return
	}

	// Validate that there are clues available to use
	if g.Clues == 0 {
		log.Fatal("The strategy of \"" + p.Strategy.Name + "\" tried to give a clue " +
			"while the team was at 0 clues.")
		return
	}

	// Validate that the clue type is sane
	if a.Clue.Type < clueTypeRank || a.Clue.Type > clueTypeColor {
		log.Fatal("The strategy of \"" + p.Strategy.Name + "\" tried to give a clue " +
			"with an invalid clue type of \"" + strconv.Itoa(a.Clue.Type) + "\".")
		return
	}

	// Validate that rank clues are valid
	if a.Clue.Type == clueTypeRank {
		valid := false
		for _, rank := range variants[g.Variant].ClueRanks {
			if rank == a.Clue.Value {
				valid = true
				break
			}
		}
		if !valid {
			log.Fatal("The strategy of \"" + p.Strategy.Name + "\" tried to give a clue " +
				"with an invalid rank of \"" + strconv.Itoa(a.Clue.Value) + "\".")
			return
		}
	}

	// Validate that the color clues are valid
	if a.Clue.Type == clueTypeColor &&
		(a.Clue.Value < 0 || a.Clue.Value > len(variants[g.Variant].ClueColors)-1) {

		log.Fatal("The strategy of \"" + p.Strategy.Name + "\" tried to give a clue " +
			"with an invalid color of \"" + strconv.Itoa(a.Clue.Value) + "\".")
		return
	}

	// Validate that the clue touches at least one card in the hand
	touchedAtLeastOneCard := false
	p2 := g.Players[a.Target]
	for _, c := range p2.Hand {
		if variantIsCardTouched(g.Variant, a.Clue, c) {
			touchedAtLeastOneCard = true
			break
		}
	}
	if !touchedAtLeastOneCard {
		log.Fatal("The strategy of \"" + p.Strategy.Name + "\" tried to give a clue " +
			"that touches no cards in the hand.")
		return
	}

	p.GiveClue(a, g)
}

func actionPlay(g *Game, p *Player, a *Action) {
	// Validate that the card is in their hand
	if !p.InHand(a.Target) {
		log.Fatal("The strategy of \"" + p.Strategy.Name + "\" tried to play a card that was " +
			"not in their hand.")
		return
	}

	c := p.RemoveCard(a.Target, g)
	p.PlayCard(g, c)
	p.DrawCard(g)
}

func actionDiscard(g *Game, p *Player, a *Action) {
	// Validate that the card is in their hand
	if !p.InHand(a.Target) {
		log.Fatal("The strategy of \"" + p.Strategy.Name + "\" tried to play a card that was " +
			"not in their hand.")
		return
	}

	// Validate that the team is not at the maximum amount of clues
	// (the client should enforce this, but do a check just in case)
	if g.Clues == maxClues {
		log.Fatal("The strategy of \"" + p.Strategy.Name + "\" tried to discard while the team " +
			"was at the maximum amount of clues.")
		return
	}

	g.Clues++
	c := p.RemoveCard(a.Target, g)
	p.DiscardCard(g, c)
	p.DrawCard(g)
}
