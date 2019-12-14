package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/cognicraft/archive"
	"github.com/cognicraft/mtg"
)

var version = "0.1"

func main() {
	name := flag.String("name", "", "Name")
	d := flag.String("data", "data.arc", "Data Archive")
	v := flag.Bool("version", false, "Version")
	flag.Parse()

	if *v {
		fmt.Printf("%s\n", version)
		os.Exit(0)
	}

	if len(flag.Args()) != 1 {
		log.Fatal(fmt.Errorf("no deck specified"))
	}

	data, err := archive.Open(*d)
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

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
	deck.Name = *name

	ext := filepath.Ext(deckFileName)
	proxyFileName := deckFileName[0:len(deckFileName)-len(ext)] + ".pdf"

	err = mtg.PDF(data, deck, proxyFileName)
	if err != nil {
		log.Fatal(err)
	}
}
