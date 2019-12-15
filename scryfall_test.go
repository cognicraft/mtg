package mtg

import (
	"testing"

	"github.com/cognicraft/archive"
)

func TestCard(t *testing.T) {
	a, _ := archive.Open(":memory:")
	s := NewScryfall(a)
	c, err := s.Card("Plains")
	if err != nil {
		t.Error(err)
	}
	if "Plains" != c.Name {
		t.Errorf("want: %s, got: %s", "Plains", c.Name)
	}
}

func TestLargeImage(t *testing.T) {
	a, _ := archive.Open(":memory:")
	s := NewScryfall(a)
	_, err := s.LargeImage("Plains")
	if err != nil {
		t.Error(err)
	}
}
