package main

import (
	"strconv"
)

type Game struct {
	Variant      string
	Players      []*Player
	Seed         int
	Deck         []*Card
	DeckIndex    int
	Stacks       []int
	Turn         int
	ActivePlayer int
	Clues        int
	Score        int
	Strikes      int
	Actions      []*Action
	EndTurn      int
}

type Action struct {
	Type   int // 0 is clue, 1 is play, 2 is discard
	Clue   Clue
	Target int
}

type Clue struct {
	Type  int // 0 is a rank clue, 1 if a clue color
	Value int
}

/*
	Initialization functions
*/

func (g *Game) InitDeck() {
	// Create the deck
	suits := make([]int, 0)
	for i := 0; i < 5; i++ {
		suits = append(suits, i) // For a normal game, the suits will be equal to {0, 1, 2, 3, 4}
	}
	for _, suit := range suits {
		ranks := []int{1, 2, 3, 4, 5}
		for _, rank := range ranks {
			// In a normal suit of Hanabi,
			// there are three 1's, two 2's, two 3's, two 4's, and one five
			var amountToAdd int
			if rank == 1 {
				amountToAdd = 3
			} else if rank == 5 {
				amountToAdd = 1
			} else {
				amountToAdd = 2
			}

			for i := 0; i < amountToAdd; i++ {
				// Add the card to the deck
				g.Deck = append(g.Deck, &Card{
					Suit: suit,
					Rank: rank,
					// We can't set the order here because the deck will be shuffled later
				})
			}
		}
	}
}

/*
	Miscellaneous functions
*/

func (g *Game) GetHandSize() int {
	numPlayers := len(g.Players)
	if numPlayers == 2 || numPlayers == 3 {
		return 5
	} else if numPlayers == 4 || numPlayers == 5 {
		return 4
	} else if numPlayers == 6 {
		return 3
	}

	log.Fatal("Failed to get the hand size for " + strconv.Itoa(numPlayers) + " players.")
	return -1
}
