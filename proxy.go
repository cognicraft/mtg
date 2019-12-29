package mtg

import (
	"bytes"
	"fmt"
	"log"

	"github.com/cognicraft/mtg/scryfall"
	"github.com/jung-kurt/gofpdf"
)

const (
	xOff       float64 = 22
	yOff       float64 = 18
	cardWidth  float64 = 63
	cardHeight float64 = 88
	labelX     float64 = 5
	labelY     float64 = 44
	labelWidth float64 = 53
	labelHight float64 = 5
)

func SimplePDF(client *scryfall.Client, deck Deck, file string) error {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()
	addCropMarks(pdf)

	pdf.SetFont("Arial", "", 8)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFillColor(255, 255, 255)

	tr := pdf.UnicodeTranslatorFromDescriptor("")

	var aux []string

	cards := deck.Cards()
	for i, card := range cards {
		col := float64(i % 4)
		row := float64((i % 8) / 4)

		sc := client.Card(card.Name)

		var name string

		switch sc.Layout {
		case scryfall.LayoutTransform:
			f := sc.Front()
			name = f.Name
			b := sc.Back()
			aux = append(aux, b.Name)
		default:
			name = sc.Name
		}
		_ = name

		x := xOff + col*cardWidth
		y := yOff + row*cardHeight
		pdf.RoundedRect(x, y, cardWidth, cardHeight, 3, "1234", "D")
		pdf.RoundedRect(x+2, y+2, cardWidth-2*2, cardHeight-2*2, 3, "1234", "D")

		pdf.MoveTo(x+2, y+2)
		pdf.CellFormat(cardWidth-2*2, 6, tr(sc.Name), "", 0, "LM", false, 0, "")
		if sc.ManaCost != "" {
			pdf.MoveTo(x+2, y+2)
			pdf.CellFormat(cardWidth-2*2, 6, tr(sc.ManaCost), "", 0, "RM", false, 0, "")
		}

		pdf.Line(x+2, y+2+6, x+cardWidth-2, y+2+6)

		imgHeight := 15.0

		pdf.Line(x+2, y+2+6, x+cardWidth-2, y+2+6+imgHeight)
		pdf.Line(x+2, y+2+6+imgHeight, x+cardWidth-2, y+2+6)
		pdf.Line(x+2, y+2+6+imgHeight, x+cardWidth-2, y+2+6+imgHeight)
		pdf.MoveTo(x+2, y+2+6+imgHeight)
		pdf.CellFormat(cardWidth-2*2, 6, tr(sc.TypeLine), "", 0, "LM", false, 0, "")
		pdf.Line(x+2, y+2+6+imgHeight+6, x+cardWidth-2, y+2+6+imgHeight+6)

		pdf.MoveTo(x+2, y+2+6+imgHeight+6+1)
		pdf.MultiCell(cardWidth-2*2, 3.8, tr(sc.OracleText), "", "LT", false)

		pdf.Line(x+2, y+cardHeight-6-2, x+cardWidth-2, y+cardHeight-6-2)
		if sc.Power != "" && sc.Toughness != "" {
			pdf.MoveTo(x+cardWidth-15, y+cardHeight-6-2-2)
			pdf.CellFormat(10, 5, fmt.Sprintf("%s / %s", sc.Power, sc.Toughness), "1", 0, "CM", true, 0, "")
		}

		if deck.Name != "" {
			pdf.MoveTo(xOff+labelX+col*cardWidth, yOff+labelY+row*cardHeight)
			pdf.CellFormat(labelWidth, labelHight, deck.Name, "", 0, "CM", true, 0, "")
		}
		if len(cards)-1 > i && i%8 == 7 {
			pdf.AddPage()
			addCropMarks(pdf)
		}
	}

	if len(aux) > 0 {
		pdf.AddPage()
		addCropMarks(pdf)
		for i, _ := range aux {
			col := float64(i % 4)
			row := float64((i % 8) / 4)
			pdf.RoundedRect(xOff+col*cardWidth, yOff+row*cardHeight, cardWidth, cardHeight, 3, "1234", "D")
			if deck.Name != "" {
				pdf.MoveTo(xOff+labelX+col*cardWidth, yOff+labelY+row*cardHeight)
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
			data, err = client.Image(f.ImageURIs["large"])

			b := sc.Back()
			if data, err := client.Image(b.ImageURIs["large"]); err == nil {
				auxImgages = append(auxImgages, img{name: b.Name, data: data})
			}
		default:
			name = sc.Name
			data, err = client.Image(sc.ImageURIs["large"])
		}
		if err != nil {
			log.Printf("ERROR: %v", err)
			continue
		}
		pdf.RegisterImageOptionsReader(name, opt, bytes.NewBuffer(data))
		pdf.ImageOptions(name, xOff+col*cardWidth, yOff+row*cardHeight, cardWidth, cardHeight, false, opt, 0, "")
		if deck.Name != "" {
			pdf.MoveTo(xOff+labelX+col*cardWidth, yOff+labelY+row*cardHeight)
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
			pdf.ImageOptions(img.name, xOff+col*cardWidth, yOff+row*cardHeight, cardWidth, cardHeight, false, opt, 0, "")
			if deck.Name != "" {
				pdf.MoveTo(xOff+labelX+col*cardWidth, yOff+labelY+row*cardHeight)
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
		pdf.Line(0, yOff+float64(i)*cardHeight, 10, yOff+float64(i)*cardHeight)
		pdf.Line(pageWidth-10, yOff+float64(i)*cardHeight, pageWidth, yOff+float64(i)*cardHeight)
	}
}

type img struct {
	name string
	data []byte
}
