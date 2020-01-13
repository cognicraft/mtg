package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cognicraft/mtg"
)

var version = "0.1"

func main() {
	n := flag.String("name", "", "Name")
	playset := flag.Bool("playset", false, "Playset?")
	v := flag.Bool("version", false, "Version")
	flag.Parse()

	if *v {
		fmt.Printf("%s\n", version)
		os.Exit(0)
	}

	if len(flag.Args()) != 1 {
		log.Fatal(fmt.Errorf("no folder specified"))
	}

	dirName := flag.Arg(0)
	outFileName := dirName + ".pdf"
	if strings.HasSuffix(dirName, "/") {
		outFileName = dirName[:len(dirName)-1] + ".pdf"
	}

	err := mtg.LayoutDirectory(*n, *playset, dirName, outFileName)
	if err != nil {
		log.Fatal(err)
	}
}
