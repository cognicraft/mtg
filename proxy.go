package mtg

import (
	"bytes"
	"fmt"
	"strings"

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

func NewProxyPrinter(client *scryfall.Client, deck Deck) *ProxyPrinter {
	return &ProxyPrinter{client: client, deck: deck}
}

type ProxyPrinter struct {
	client *scryfall.Client
	deck   Deck
}

func (p *ProxyPrinter) WriteImageProxies(file string) error {
	pdf := gofpdf.New("L", "mm", "A4", "")

	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(255, 255, 255)

	opt := gofpdf.ImageOptions{
		ImageType:             "jpg",
		AllowNegativePosition: true,
	}

	writeSection := func(cards []card) {
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
			if p.deck.Name != "" {
				pdf.MoveTo(xOff+labelX+col*cardWidth, yOff+labelY+row*cardHeight)
				pdf.CellFormat(labelWidth, labelHight, p.deck.Name, "", 0, "CM", true, 0, "")
			}
			if len(cards)-1 > i && i%8 == 7 {
				pdf.AddPage()
				addCropMarks(pdf)
			}
		}
	}

	deck := p.collectDeck()
	writeSection(deck.frontCards)
	writeSection(deck.backCards)
	writeSection(deck.tokens)

	return pdf.OutputFileAndClose(file)
}

func (p *ProxyPrinter) WriteTextProxies(file string) error {
	pdf := gofpdf.New("L", "mm", "A4", "")

	pdf.SetFont("Arial", "", 8)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFillColor(255, 255, 255)

	rep := strings.NewReplacer("−", "-")
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	writeSection := func(cards []card) {
		pdf.AddPage()
		addCropMarks(pdf)
		for i, card := range cards {
			col := float64(i % 4)
			row := float64((i % 8) / 4)
			x := xOff + col*cardWidth
			y := yOff + row*cardHeight
			pdf.RoundedRect(x, y, cardWidth, cardHeight, 3, "1234", "D")
			pdf.RoundedRect(x+2, y+2, cardWidth-2*2, cardHeight-2*2, 3, "1234", "D")

			pdf.MoveTo(x+2, y+2)
			pdf.CellFormat(cardWidth-2*2, 6, tr(card.Name), "", 0, "LM", false, 0, "")
			if card.ManaCost != "" {
				pdf.MoveTo(x+2, y+2)
				pdf.CellFormat(cardWidth-2*2, 6, tr(card.ManaCost), "", 0, "RM", false, 0, "")
			}

			pdf.Line(x+2, y+2+6, x+cardWidth-2, y+2+6)

			imgHeight := 15.0

			pdf.Line(x+2, y+2+6, x+cardWidth-2, y+2+6+imgHeight)
			pdf.Line(x+2, y+2+6+imgHeight, x+cardWidth-2, y+2+6)
			pdf.Line(x+2, y+2+6+imgHeight, x+cardWidth-2, y+2+6+imgHeight)
			pdf.MoveTo(x+2, y+2+6+imgHeight)
			pdf.CellFormat(cardWidth-2*2, 6, tr(card.TypeLine), "", 0, "LM", false, 0, "")
			pdf.Line(x+2, y+2+6+imgHeight+6, x+cardWidth-2, y+2+6+imgHeight+6)

			pdf.MoveTo(x+2, y+2+6+imgHeight+6+1)
			pdf.MultiCell(cardWidth-2*2, 3.8, tr(rep.Replace(card.OracleText)), "", "LT", false)

			pdf.Line(x+2, y+cardHeight-6-2, x+cardWidth-2, y+cardHeight-6-2)
			if card.Power != "" && card.Toughness != "" {
				pdf.MoveTo(x+cardWidth-15, y+cardHeight-6-2-2)
				pdf.CellFormat(10, 5, fmt.Sprintf("%s / %s", card.Power, card.Toughness), "1", 0, "CM", true, 0, "")
			}
			if card.Loyalty != "" {
				pdf.MoveTo(x+cardWidth-15, y+cardHeight-6-2-2)
				pdf.CellFormat(10, 5, fmt.Sprintf("%s", card.Loyalty), "1", 0, "CM", true, 0, "")
			}

			if p.deck.Name != "" {
				pdf.MoveTo(xOff+labelX+col*cardWidth, yOff+labelY+row*cardHeight)
				pdf.CellFormat(labelWidth, labelHight, p.deck.Name, "", 0, "CM", true, 0, "")
			}
			if len(cards)-1 > i && i%8 == 7 {
				pdf.AddPage()
				addCropMarks(pdf)
			}
		}
	}

	deck := p.collectDeck()
	writeSection(deck.frontCards)
	writeSection(deck.backCards)
	writeSection(deck.tokens)

	return pdf.OutputFileAndClose(file)
}

func (p *ProxyPrinter) collectDeck() deck {
	cardFromCard := func(sc *scryfall.Card) card {
		c := card{
			Name:       sc.Name,
			ManaCost:   sc.ManaCost,
			TypeLine:   sc.TypeLine,
			OracleText: sc.OracleText,
			Power:      sc.Power,
			Toughness:  sc.Toughness,
			Loyalty:    sc.Loyalty,
		}
		if data, err := p.client.ImageByURL(sc.ImageURIs["large"]); err == nil {
			c.ImageData = data
		}
		return c
	}

	cardFromFace := func(sc *scryfall.CardFace) card {
		c := card{
			Name:       sc.Name,
			ManaCost:   sc.ManaCost,
			TypeLine:   sc.TypeLine,
			OracleText: sc.OracleText,
			Power:      sc.Power,
			Toughness:  sc.Toughness,
			Loyalty:    sc.Loyalty,
		}
		if data, err := p.client.ImageByURL(sc.ImageURIs["large"]); err == nil {
			c.ImageData = data
		}
		return c
	}
	d := deck{}

	cards := p.deck.Cards()
	for _, card := range cards {

		sc := p.client.CardByName(card.Name)
		if sc == nil {
			continue
		}

		switch sc.Layout {
		case scryfall.LayoutTransform:
			d.frontCards = append(d.frontCards, cardFromFace(sc.Front()))
			d.backCards = append(d.backCards, cardFromFace(sc.Back()))
		default:
			d.frontCards = append(d.frontCards, cardFromCard(sc))
		}

		if len(sc.AllParts) > 0 {
			for _, part := range sc.AllParts {
				if part.Component == "token" {
					if tc := p.client.CardByURL(part.URI); tc != nil {
						d.tokens = append(d.tokens, cardFromCard(tc))
					}
				}
			}
		}
	}
	return d
}

type deck struct {
	frontCards []card
	backCards  []card
	tokens     []card
}

type card struct {
	Name       string
	ManaCost   string
	TypeLine   string
	OracleText string
	Power      string
	Toughness  string
	Loyalty    string
	ImageData  []byte
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
