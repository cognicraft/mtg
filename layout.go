package mtg

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/jung-kurt/gofpdf"
)

func LayoutFolder(inDir string, outFile string) error {
	deckName := ""
	pdf := gofpdf.New("L", "mm", "A4", "")

	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(255, 255, 255)

	opt := gofpdf.ImageOptions{
		ImageType:             "jpg",
		AllowNegativePosition: true,
	}

	writeSection := func(cards []Card) {
		if len(cards) == 0 {
			// empty section
			return
		}
		pdf.AddPage()
		addCropMarks(pdf)
		for i, card := range cards {
			col := float64(i % 4)
			row := float64((i % 8) / 4)

			pdf.RegisterImageOptionsReader(card.Name, opt, bytes.NewBuffer(card.ImageData))
			pdf.ImageOptions(card.Name, xOff+col*cardWidth, yOff+row*cardHeight, cardWidth, cardHeight, false, opt, 0, "")
			if deckName != "" {
				pdf.MoveTo(xOff+labelX+col*cardWidth, yOff+labelY+row*cardHeight)
				pdf.CellFormat(labelWidth, labelHight, deckName, "", 0, "CM", true, 0, "")
			}
			if len(cards)-1 > i && i%8 == 7 {
				pdf.AddPage()
				addCropMarks(pdf)
			}
		}
	}

	deck := Deck{}
	var main Section

	files, err := ioutil.ReadDir(inDir)
	if err != nil {
		return err
	}
	for _, file := range files {
		ext := filepath.Ext(file.Name())
		switch ext {
		case ".jpg":
			c := Card{Name: file.Name()}
			if data, err := ioutil.ReadFile(file.Name()); err != nil {
				c.ImageData = data
			}
			main.Cards = append(main.Cards, c)
		}
		fmt.Println(file.Name())
	}

	for _, s := range deck.Sections {
		writeSection(s.Cards)
	}

	return pdf.OutputFileAndClose(outFile)
}
