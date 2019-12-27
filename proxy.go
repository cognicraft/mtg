package mtg

import (
	"bytes"
	"log"

	"github.com/cognicraft/mtg/scryfall"
	"github.com/jung-kurt/gofpdf"
)

const (
	xOff       float64 = 22
	yOff       float64 = 18
	cardWidth  float64 = 63
	cardHight  float64 = 88
	labelX     float64 = 5
	labelY     float64 = 44
	labelWidth float64 = 53
	labelHight float64 = 5
)

func PDF(client *scryfall.Client, deck Deck, file string) error {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()
	addCropMarks(pdf)

	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(255, 255, 255)

	opt := gofpdf.ImageOptions{
		ImageType:             "jpg",
		AllowNegativePosition: true,
	}

	var auxImgages []img

	cards := deck.Cards()
	for i, card := range cards {
		col := float64(i % 4)
		row := float64((i % 8) / 4)

		sc := client.Card(card.Name)

		var name string
		var data []byte
		var err error

		switch sc.Layout {
		case scryfall.LayoutTransform:
			f := sc.Front()
			name = f.Name
			data, err = f.Image("large")

			b := sc.Back()
			if data, err := b.Image("large"); err == nil {
				auxImgages = append(auxImgages, img{name: b.Name, data: data})
			}
		default:
			name = sc.Name
			data, err = sc.Image("large")
		}
		if err != nil {
			log.Printf("ERROR: %v", err)
			continue
		}
		pdf.RegisterImageOptionsReader(name, opt, bytes.NewBuffer(data))
		pdf.ImageOptions(name, xOff+col*cardWidth, yOff+row*cardHight, cardWidth, cardHight, false, opt, 0, "")
		if deck.Name != "" {
			pdf.MoveTo(xOff+labelX+col*cardWidth, yOff+labelY+row*cardHight)
			pdf.CellFormat(labelWidth, labelHight, deck.Name, "", 0, "CM", true, 0, "")
		}
		if len(cards)-1 > i && i%8 == 7 {
			pdf.AddPage()
			addCropMarks(pdf)
		}
	}

	if len(auxImgages) > 0 {
		pdf.AddPage()
		addCropMarks(pdf)
		for i, img := range auxImgages {
			col := float64(i % 4)
			row := float64((i % 8) / 4)

			pdf.RegisterImageOptionsReader(img.name, opt, bytes.NewBuffer(img.data))
			pdf.ImageOptions(img.name, xOff+col*cardWidth, yOff+row*cardHight, cardWidth, cardHight, false, opt, 0, "")
			if deck.Name != "" {
				pdf.MoveTo(xOff+labelX+col*cardWidth, yOff+labelY+row*cardHight)
				pdf.CellFormat(labelWidth, labelHight, deck.Name, "", 0, "CM", true, 0, "")
			}
			if len(cards)-1 > i && i%8 == 7 {
				pdf.AddPage()
				addCropMarks(pdf)
			}
		}
	}

	return pdf.OutputFileAndClose(file)
}

func addCropMarks(pdf *gofpdf.Fpdf) {
	pageWidth, pageHeight := pdf.GetPageSize()
	for i := 0; i <= 4; i++ {
		pdf.Line(xOff+float64(i)*cardWidth, 0, xOff+float64(i)*cardWidth, 10)
		pdf.Line(xOff+float64(i)*cardWidth, pageHeight-10, xOff+float64(i)*cardWidth, pageHeight)
	}
	for i := 0; i <= 2; i++ {
		pdf.Line(0, yOff+float64(i)*cardHight, 10, yOff+float64(i)*cardHight)
		pdf.Line(pageWidth-10, yOff+float64(i)*cardHight, pageWidth, yOff+float64(i)*cardHight)
	}
}

type img struct {
	name string
	data []byte
}
