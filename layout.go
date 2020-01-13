package mtg

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/jung-kurt/gofpdf"
)

func LayoutDirectory(deckName string, numberOfCopiesPerCard int, inDir string, outFile string) error {
	pdf := gofpdf.New("L", "mm", "A4", "")

	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(255, 255, 255)

	optFor := func(fileName string) gofpdf.ImageOptions {
		ext := filepath.Ext(fileName)
		switch ext {
		case ".jpg", ".jpeg", ".png", ".gif":
			return gofpdf.ImageOptions{
				ImageType:             ext[1:],
				AllowNegativePosition: true,
			}
		}
		return gofpdf.ImageOptions{}
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

			opt := optFor(card.Name)
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

	var main Section
	files, err := ioutil.ReadDir(inDir)
	if err != nil {
		return err
	}
	for _, file := range files {
		ext := filepath.Ext(file.Name())
		switch ext {
		case ".jpg", ".jpeg", ".png", ".gif":
			c := Card{Name: file.Name()}
			if data, err := ioutil.ReadFile(filepath.Join(inDir, file.Name())); err == nil {
				c.ImageData = data
			} else {
				fmt.Printf("ERROR: %v\n", err)
			}
			for i := 0; i < numberOfCopiesPerCard; i++ {
				main.Cards = append(main.Cards, c)
			}
		}
	}
	if len(main.Cards) > 0 {
		writeSection(main.Cards)
	}

	return pdf.OutputFileAndClose(outFile)
}
