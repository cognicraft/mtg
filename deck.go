package mtg

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

type Deck struct {
	Name  string
	Cards []Card
}

type Card struct {
	Name string
}

func ParseDeck(in io.Reader) (Deck, error) {
	deck := Deck{}
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "//") {
			continue
		}
		fs := strings.Fields(line)
		if len(fs) < 2 {
			continue
		}
		c, err := strconv.ParseInt(fs[0], 10, 64)
		if err != nil {
			continue
		}
		name := strings.Join(fs[1:], " ")
		n := int(c)
		for i := 0; i < n; i++ {
			deck.Cards = append(deck.Cards, Card{Name: name})
		}
	}
	return deck, nil
}
