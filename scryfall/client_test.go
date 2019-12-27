package scryfall

import (
	"testing"
)

func TestCard(t *testing.T) {
	s, err := New(Debug)
	if err != nil {
		t.Error(err)
	}
	c := s.Card("Plains")
	if "Plains" != c.Name {
		t.Errorf("want: %s, got: %s", "Plains", c.Name)
	}
	if !c.IsLegalIn("pioneer") {
		t.Errorf("expected Plains to be legal in Pioneer")
	}
	c.Image("large")

	// t.Fail()
}

func TestNicol(t *testing.T) {
	c, _ := New(Debug)
	card := c.Card("Nicol Bolas, the Ravager")
	t.Logf("%v", card.Front())
	t.Fail()
}
