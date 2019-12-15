package mtg

import (
	"bytes"
	"log"

	"github.com/cognicraft/archive"
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

func PDF(data *archive.Archive, deck Deck, file string) error {
	s := NewScryfall(data)
	pdf := gofpdf.New("L", "mm", "A4", "")

	pdf.AddPage()
	addCropMarks(pdf)

	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(255, 255, 255)

	opt := gofpdf.ImageOptions{
		ImageType:             "jpg",
		AllowNegativePosition: true,
	}

	for i, c := range deck.Cards {
		col := float64(i % 4)
		row := float64((i % 8) / 4)

		bs, err := s.LargeImage(c.Name)
		if err != nil {
			log.Printf("ERROR: %v", err)
			continue
		}
		pdf.RegisterImageOptionsReader(c.Name, opt, bytes.NewBuffer(bs))
		pdf.ImageOptions(c.Name, xOff+col*cardWidth, yOff+row*cardHight, cardWidth, cardHight, false, opt, 0, "")
		if deck.Name != "" {
			pdf.MoveTo(xOff+labelX+col*cardWidth, yOff+labelY+row*cardHight)
			pdf.CellFormat(labelWidth, labelHight, deck.Name, "", 0, "CM", true, 0, "")
		}
		if len(deck.Cards)-1 > i && i%8 == 7 {
			pdf.AddPage()
			addCropMarks(pdf)
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
