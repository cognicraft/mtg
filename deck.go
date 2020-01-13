package mtg

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

type Deck struct {
	Name     string
	Sections []Section
}

func (d Deck) Cards() []Card {
	var cards []Card
	for _, s := range d.Sections {
		for _, c := range s.Cards {
			cards = append(cards, c)
		}
	}
	return cards
}

type Section struct {
	Name  string
	Cards Cards
}

type Cards []Card

func (cs Cards) Contains(accept func(Card) bool) bool {
	for _, c := range cs {
		if accept(c) {
			return true
		}
	}
	return false
}

func CardByName(name string) func(Card) bool {
	return func(c Card) bool {
		return c.Name == name
	}
}

type Card struct {
	Name       string
	ManaCost   string
	TypeLine   string
	OracleText string
	Power      string
	Toughness  string
	Loyalty    string
	ImageData  []byte
	Version    *Version
}

type Version struct {
	Set             string
	CollectorNumber string
}

func ParseDeck(in io.Reader) (Deck, error) {
	deck := Deck{}
	currentSection := Section{Name: "Main"}
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "//") {
			continue
		}
		if i := strings.Index(line, "#"); i >= 0 {
			line = line[:i]
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		card := Card{}
		if i := strings.Index(line, "["); i >= 0 {
			if o := strings.Index(line, "]"); o > i {
				version := line[i+1 : o]
				vps := strings.Split(version, ":")
				switch len(vps) {
				case 1:
					card.Version = &Version{Set: vps[0]}
				case 2:
					card.Version = &Version{Set: vps[0], CollectorNumber: vps[1]}
				}
				line = line[:i] + line[o+1:]
			}
		}
		fs := strings.Fields(line)
		if len(fs) < 2 {
			continue
		}
		c, err := strconv.ParseInt(fs[0], 10, 64)
		if err != nil {
			continue
		}
		card.Name = strings.Join(fs[1:], " ")
		n := int(c)
		for i := 0; i < n; i++ {
			currentSection.Cards = append(currentSection.Cards, card)
		}
	}
	deck.Sections = append(deck.Sections, currentSection)
	return deck, nil
}
