package main

import (
	"math/rand"
	"os"
	"strconv"

	"github.com/op/go-logging"
)

const (
	numPlayers   = 3
	stratToUse   = "Hyphen-ated"
	variantToUse = "No Variant"
)

var (
	names = []string{"Alice", "Bob", "Cathy", "Donald", "Emily"}
	log   *logging.Logger
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
		Variant:       variantToUse,
		Players:       make([]*Player, 0),
		PossibleCards: make(map[string]int),
		Stacks:        make([]int, 0),
		DiscardPile:   make([]*Card, 0),
		Clues:         maxClues,
		Actions:       make([]*Action, 0),
		EndTurn:       -1,
	}

	g.InitDeck()
	g.InitStacks()
	rand.Seed(int64(g.Seed)) // Seed the random number generator with the game seed
	g.Shuffle()
	g.InitPlayers()
	g.DealStartingHands()

	// Allow the strategies to "see" the opening hands
	for i, p := range g.Players {
		p.Strategy.Start(p.Strategy, g, i)
	}

	// Play the game until it ends
	for {
		// Query the strategy to see what kind of move that the player will do
		p := g.Players[g.ActivePlayer]
		a := p.Strategy.GetAction(p.Strategy, g)
		if a == nil {
			log.Fatal("The strategy of \"" + p.Strategy.Name + "\" returned a nil action.")
		}

		// Allow the strategies to "see" the upcoming action
		for _, p := range g.Players {
			p.Strategy.ActionAnnounced(p.Strategy, g, a)
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

		// Allow the strategies to "see" the game state after the action is completed
		for _, p := range g.Players {
			p.Strategy.ActionHappened(p.Strategy, g, a)
		}

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
		} else {
			log.Info("It is now " + g.Players[g.ActivePlayer].Name + "'s turn.")
		}

		if g.EndCondition > endConditionInProgress {
			break
		}
	}

	// Provide a JSON export of the game that can be imported into Hanabi Live
	g.Export()
}
