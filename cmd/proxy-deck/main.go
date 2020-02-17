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
	f := flag.String("format", "image", "Format")
	withTokens := flag.Bool("with-tokens", false, "With tokens?")
	onlyTokens := flag.Bool("only-tokens", false, "Print only associated tokens")
	numberOfTokens := flag.Int("number-of-tokens", 4, "The number of each token to print.")
	debug := flag.Bool("debug", false, "Debug?")
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

	var scOpts []scryfall.ClientOption
	scOpts = append(scOpts, scryfall.Cache(cache))
	if *debug {
		scOpts = append(scOpts, scryfall.Debug())
	}

	scry, err := scryfall.New(scOpts...)
	if err != nil {
		log.Fatal(err)
	}

	var opts []mtg.PrinterOption
	opts = append(opts, mtg.NumberOfTokens(*numberOfTokens))
	if *withTokens {
		opts = append(opts, mtg.PrintTokens())
	}
	if *onlyTokens {
		opts = append(opts, mtg.PrintOnlyTokens())
	}

	switch *f {
	case "text":
		err = mtg.NewProxyPrinter(scry, deck, opts...).WriteTextProxiesToFile(proxyFileName)
	default:
		err = mtg.NewProxyPrinter(scry, deck, opts...).WriteImageProxiesToFile(proxyFileName)
	}
	if err != nil {
		log.Fatal(err)
	}
}
