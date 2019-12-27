package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/cognicraft/archive"
	"github.com/cognicraft/mtg"
	"github.com/cognicraft/mtg/scryfall"
)

var version = "0.1"

func main() {
	n := flag.String("name", "", "Name")
	c := flag.String("cache", "cache.arc", "Cache")
	v := flag.Bool("version", false, "Version")
	flag.Parse()

	if *v {
		fmt.Printf("%s\n", version)
		os.Exit(0)
	}

	if len(flag.Args()) != 1 {
		log.Fatal(fmt.Errorf("no deck specified"))
	}

	cache, err := archive.Open(*c)
	if err != nil {
		log.Fatal(err)
	}
	defer cache.Close()

	deckFileName := flag.Arg(0)
	deckFile, err := os.Open(deckFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer deckFile.Close()

	deck, err := mtg.ParseDeck(deckFile)
	if err != nil {
		log.Fatal(err)
	}
	deck.Name = *n

	ext := filepath.Ext(deckFileName)
	proxyFileName := deckFileName[0:len(deckFileName)-len(ext)] + ".pdf"

	scry, err := scryfall.New(
		scryfall.Cache(cache),
		scryfall.Debug,
	)
	if err != nil {
		log.Fatal(err)
	}

	err = mtg.PDF(scry, deck, proxyFileName)
	if err != nil {
		log.Fatal(err)
	}
}
