package main

type Variant struct {
	Suits      []*Suit
	Ranks      []int
	ClueColors []string
	ClueRanks  []int
}

type Suit struct {
	Name       string
	Index      int
	ClueColors []string
	OneOfEach  bool
}

func (s *Suit) ColorValue(g *Game) int {
	if len(s.ClueColors) > 1 {
		log.Fatal("Dual-color variants are not currently implemented.")
	}

	for i, color := range variants[g.Variant].ClueColors {
		if color == s.ClueColors[0] {
			return i
		}
	}

	return -1
}

var (
	variants map[string]*Variant
)

func variantsInit() {
	variants = make(map[string]*Variant)

	suits := make([]*Suit, 0)
	colors := []string{"Blue", "Green", "Yellow", "Red", "Purple"}
	for i, color := range colors {
		suits = append(suits, &Suit{
			Name:       color,
			Index:      i,
			ClueColors: []string{color},
		})
	}
	variants["No Variant"] = &Variant{
		Suits:      suits,
		Ranks:      []int{1, 2, 3, 4, 5},
		ClueColors: colors,
		ClueRanks:  []int{1, 2, 3, 4, 5},
	}
}

// variantIsCardTouched returns true if a clue will touch a particular suit
// For example, a yellow clue will not touch a green card in a normal game,
// but it will the "Dual-Color" variant
func variantIsCardTouched(g *Game, clue *Clue, card *Card) bool {
	if clue.Type == clueTypeRank {
		return card.Rank == clue.Value
	} else if clue.Type == clueTypeColor {
		color := variants[g.Variant].ClueColors[clue.Value]
		colors := variants[g.Variant].Suits[card.Suit.Index].ClueColors
		return stringInSlice(color, colors)
	}

	return false
}
