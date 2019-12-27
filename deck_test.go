package mtg

import (
	"strings"
	"testing"
)

func TestParseDeck(t *testing.T) {
	raw := `
//Main
1 [DOM:205] Slimefoot, the Stowaway # !Commander
`

	d, err := ParseDeck(strings.NewReader(raw))
	if err != nil {
		t.Error(err)
	}
	cs := d.Cards()
	if len(cs) != 1 {
		t.Fatalf("want: %d, got: %d", 1, len(cs))
	}
	c := cs[0]
	if "Slimefoot, the Stowaway" != c.Name {
		t.Errorf("want: %q, got: %q", "Slimefoot, the Stowaway", cs[0].Name)
	}
	if c.Version == nil {
		t.Fatal("expected version")
	}
	want := Version{Set: "DOM", CollectorNumber: "205"}
	if want != *c.Version {
		t.Errorf("want: %v, got: %v", want, *c.Version)
	}

}
