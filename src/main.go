package main

import (
	"math/rand"
	"os"
	"strconv"

	"github.com/op/go-logging"
)

const (
	numPlayers   = 3
	stratToUse   = "Dumb"
	variantToUse = "No Variant"
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
		Variant:     variantToUse,
		Players:     make([]*Player, 0),
		Stacks:      make([]int, 0),
		DiscardPile: make([]*Card, 0),
		Clues:       maxClues,
		Actions:     make([]*Action, 0),
		EndTurn:     -1,
	}

	g.InitDeck()
	g.InitStacks()
	rand.Seed(int64(g.Seed)) // Seed the random number generator with the game seed
	g.Shuffle()
	g.InitPlayers()
	g.DealStartingHands()

	// Allow the strategies to "see" the opening hands
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

	g.Export()
}
