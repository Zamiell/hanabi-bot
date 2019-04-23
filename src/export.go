package main

import (
	"encoding/json"

	"github.com/atotto/clipboard"
)

type JSONGame struct {
	Actions     []*Action  `json:"actions"`
	Deck        []JSONCard `json:"deck"`
	FirstPlayer int        `json:"firstPlayer"`
	Notes       [][]string `json:"notes"`
	Players     []string   `json:"players"`
	Variant     string     `json:"variant"`
}
type JSONCard struct {
	Suit int `json:"suit"`
	Rank int `json:"rank"`
}

// Export creates a JSON export of the game (in the Hanabi Live format)
// https://github.com/Zamiell/hanabi-live/blob/master/misc/example_game_with_comments.json
func (g *Game) Export() {
	// Create the deck
	deck := make([]JSONCard, 0)
	for _, c := range g.Deck {
		deck = append(deck, JSONCard{
			Suit: c.Suit.Index,
			Rank: c.Rank,
		})
	}

	// Create the notes
	notes := make([][]string, 0)
	for _, p := range g.Players {
		notes = append(notes, p.Notes)
	}

	// Compile the player names
	players := make([]string, 0)
	for _, p := range g.Players {
		players = append(players, p.Name)
	}

	// Instantiate the object
	exportObject := &JSONGame{
		Actions:     g.Actions,
		Deck:        deck,
		FirstPlayer: g.FirstPlayer,
		Notes:       notes,
		Players:     players,
		Variant:     g.Variant,
	}

	// Convert it to JSON
	var exportJSON []byte
	if v, err := json.Marshal(exportObject); err != nil {
		log.Fatal("Failed to convert the exported object to JSON:", err)
	} else {
		exportJSON = v
	}

	// Copy the JSON to the clipboard
	if err := clipboard.WriteAll(string(exportJSON)); err != nil {
		log.Fatal("Failed to write to the clipboard:", err)
		return
	}
}
