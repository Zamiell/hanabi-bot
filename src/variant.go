package main

import (
	"strings"
)

type Variant struct {
	Suits      []*Suit
	ClueColors []string
	ClueRanks  []int
}

type Suit struct {
	Name       string
	ClueColors []string
	OneOfEach  bool
}

var (
	variants map[string]*Variant
)

func variantsInit() {
	variants = make(map[string]*Variant)

	suits := make([]*Suit, 0)
	colors := []string{"Blue", "Green", "Yellow", "Red", "Purple"}
	for _, color := range colors {
		suits = append(suits, &Suit{
			Name:       color,
			ClueColors: []string{color},
		})
	}
	variants["No Variant"] = &Variant{
		Suits:      suits,
		ClueColors: colors,
		ClueRanks:  []int{1, 2, 3, 4, 5},
	}
}

// variantIsCardTouched returns true if a clue will touch a particular suit
// For example, a yellow clue will not touch a green card in a normal game,
// but it will the "Dual-Color" variant
func variantIsCardTouched(variant string, clue Clue, card *Card) bool {
	if strings.HasPrefix(variant, "Totally Mute") {
		return false
	}

	if clue.Type == clueTypeRank {
		return card.Rank == clue.Value
	} else if clue.Type == clueTypeColor {
		color := variants[variant].ClueColors[clue.Value]
		colors := variants[variant].Suits[card.Suit].ClueColors
		return stringInSlice(color, colors)
	}

	return false
}
