package scryfall

import (
	"testing"
)

func TestCard(t *testing.T) {
	s, err := New(Debug)
	if err != nil {
		t.Error(err)
	}
	c := s.CardByName("Plains")
	if "Plains" != c.Name {
		t.Errorf("want: %s, got: %s", "Plains", c.Name)
	}
	if !c.IsLegalIn("pioneer") {
		t.Errorf("expected Plains to be legal in Pioneer")
	}
	s.ImageByURL(c.ImageURIs["large"])

	// t.Fail()
}

func TestNicol(t *testing.T) {
	c, _ := New(Debug)
	card := c.CardByName("Nicol Bolas, the Ravager")
	t.Logf("%v", card.Front())
}
