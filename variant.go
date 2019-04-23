package main

type Variant struct {
	Suits []*Suit
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
		Suits: suits,
	}
}
