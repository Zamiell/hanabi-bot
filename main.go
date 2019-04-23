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
		Variant: "No Variant",
		Players: make([]*Player, 0),
		Stacks:  make([]int, 0),
		Clues:   8,
		Actions: make([]*Action, 0),
		EndTurn: -1,
	}

	// Initialize the players
	names := []string{"Alice", "Bob"}
	for i := 0; i < 2; i++ {
		p := &Player{
			Name: names[i],
			Hand: make([]*Card, 0),
		}
		g.Players = append(g.Players, p)
	}

	g.InitDeck()

	// Initialize the stacks
	g.Stacks = make([]int, 0)
	for i := 0; i < 5; i++ {
		g.Stacks = append(g.Stacks, 0)
	}

	// Seed the random number generator with the game seed
	rand.Seed(int64(g.Seed))

	// Shuffle the deck
	// From: https://stackoverflow.com/questions/12264789/shuffle-array-in-go
	for i := range g.Deck {
		j := rand.Intn(i + 1)
		g.Deck[i], g.Deck[j] = g.Deck[j], g.Deck[i]
	}

	// Mark the order of all of the cards in the deck
	for i, c := range g.Deck {
		c.Order = i
	}

	// Log the deal (so that it can be distributed to others if necessary)
	log.Info("--------------------------------------------------")
	log.Info("Deal for seed: " + strconv.Itoa(g.Seed) + " (from top to bottom)")
	log.Info("(cards are dealt to a player until their hand fills up before moving on to the next one)")
	for i, c := range g.Deck {
		log.Info(strconv.Itoa(i+1) + ") " + c.Name(g))
	}
	log.Info("--------------------------------------------------")

	// Get a random player to start first (based on the game seed)
	g.ActivePlayer = rand.Intn(len(g.Players))

	// Shuffle the order of the players
	// (otherwise, the seat order would always correspond to the order that
	// the players joined the game in)
	// From: https://stackoverflow.com/questions/12264789/shuffle-array-in-go
	for i := range g.Players {
		j := rand.Intn(i + 1)
		g.Players[i], g.Players[j] = g.Players[j], g.Players[i]
	}

	// Set the player indexes
	for i, p := range g.Players {
		p.Index = i
	}

	// Deal the cards
	handSize := g.GetHandSize()
	for _, p := range g.Players {
		for i := 0; i < handSize; i++ {
			p.DrawCard(g)
		}
	}

}
