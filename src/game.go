package main

import (
	"math/rand"
	"strconv"
)

type Game struct {
	Variant       string
	Players       []*Player
	Seed          int
	Deck          []*Card
	PossibleCards map[string]int
	DeckIndex     int
	Stacks        []int
	DiscardPile   []*Card
	Turn          int
	FirstPlayer   int
	ActivePlayer  int
	Clues         int
	Score         int
	Strikes       int
	Actions       []*Action
	EndCondition  int
	EndTurn       int
}

/*
	Initialization functions
*/

func (g *Game) InitDeck() {
	// Create the deck
	for _, suit := range variants[g.Variant].Suits {
		for _, rank := range variants[g.Variant].Ranks {
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
				c := &Card{
					Suit: suit,
					Rank: rank,
					// We can't set the order here because the deck will be shuffled later
					Holder:        -1,
					Slot:          -1,
					Clues:         make([]*CardClue, 0),
					PossibleSuits: append(variants[g.Variant].Suits[:0:0], variants[g.Variant].Suits...),
					PossibleRanks: append(variants[g.Variant].Ranks[:0:0], variants[g.Variant].Ranks...),
					PossibleCards: make(map[string]int),
				}

				g.Deck = append(g.Deck, c)

				// Add the possibility
				mapIndex := suit.Name + strconv.Itoa(rank)
				g.PossibleCards[mapIndex]++
			}
		}
	}

	// Copy all of the possibilities into every card
	for _, c := range g.Deck {
		for k, v := range g.PossibleCards {
			c.PossibleCards[k] = v
		}
	}
}

func (g *Game) InitStacks() {
	g.Stacks = make([]int, 0)
	for i := 0; i < len(variants[g.Variant].Suits); i++ {
		g.Stacks = append(g.Stacks, 0)
	}
}

func (g *Game) Shuffle() {
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

	/*
		// Log the deal (so that it can be distributed to others if necessary)
		log.Info("--------------------------------------------------")
		log.Info("Deal for seed: " + strconv.Itoa(g.Seed) + " (from top to bottom)")
		log.Info("(cards are dealt to a player until their hand fills up before moving on to the next one)")
		for i, c := range g.Deck {
			log.Info(strconv.Itoa(i+1) + ") " + c.Name(g))
		}
		log.Info("--------------------------------------------------")
	*/
}

func (g *Game) InitPlayers() {
	for i := 0; i < numPlayers; i++ {
		notes := make([]string, 0)
		for j := 0; j < len(g.Deck); j++ {
			notes = append(notes, "")
		}

		p := &Player{
			Name:     names[i],
			Hand:     make([]*Card, 0),
			Notes:    notes,
			Strategy: strategies[stratToUse](),
		}
		g.Players = append(g.Players, p)
	}

	// Get a random player to start first (based on the game seed)
	g.FirstPlayer = rand.Intn(len(g.Players))
	g.ActivePlayer = g.FirstPlayer

	// Shuffle the order of the players
	// (otherwise, the seat order would always correspond to the order that
	// the players joined the game in)
	// From: https://stackoverflow.com/questions/12264789/shuffle-array-in-go
	for i := range g.Players {
		j := rand.Intn(i + 1)
		g.Players[i], g.Players[j] = g.Players[j], g.Players[i]
	}

	// Store the player indexes for later
	for i, p := range g.Players {
		p.Index = i
	}

	log.Info("The seating of the players is as follows:")
	for i, p := range g.Players {
		str := strconv.Itoa(i+1) + ") " + p.Name
		log.Info(str)
	}
	log.Info("----------------------------------------")
}

func (g *Game) DealStartingHands() {
	log.Info("Performing the initial deal:")
	handSize := g.GetHandSize()
	for _, p := range g.Players {
		for i := 0; i < handSize; i++ {
			p.DrawCard(g)
		}
	}
	log.Info("----------------------------------------")
	log.Info(g.Players[g.ActivePlayer].Name + " goes first.")
}

/*
	Main functions
*/

func (g *Game) CheckEnd() bool {
	// Check for 3 strikes
	if g.Strikes == 3 {
		g.EndCondition = endConditionStrikeout
		return true
	}

	// Check to see if the final go-around has completed
	// (which is initiated after the last card is played from the deck)
	if g.Turn == g.EndTurn {
		g.EndCondition = endConditionNormal
		return true
	}

	return false
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

// GetSpecificCardNum returns the total cards in the deck of the specified suit and rank
// as well as how many of those that have been already discarded
func (g *Game) GetSpecificCardNum(suit *Suit, rank int) (int, int) {
	total := 0
	discarded := 0
	for _, c := range g.Deck {
		if c.Suit == suit && c.Rank == rank {
			total++
			if c.Discarded {
				discarded++
			}
		}
	}

	return total, discarded
}
