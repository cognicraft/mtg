package mtg

import (
	"bytes"
	"fmt"
	"io"
	"os"
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

type PrinterOption func(*ProxyPrinter) error

func PrintTokens() PrinterOption {
	return func(p *ProxyPrinter) error {
		p.printTokens = true
		return nil
	}
}

func PrintOnlyTokens() PrinterOption {
	return func(p *ProxyPrinter) error {
		p.printFrontFaces = false
		p.printBackFaces = false
		p.printTokens = true
		return nil
	}
}

func NumberOfTokens(n int) PrinterOption {
	return func(p *ProxyPrinter) error {
		p.numberOfTokens = n
		return nil
	}
}

func Language(lang scryfall.Lang) PrinterOption {
	return func(p *ProxyPrinter) error {
		p.lang = lang
		return nil
	}
}

func NewProxyPrinter(client *scryfall.Client, deck Deck, opts ...PrinterOption) *ProxyPrinter {
	p := &ProxyPrinter{
		client:          client,
		lang:            scryfall.LangEnglish,
		deck:            deck,
		printFrontFaces: true,
		printBackFaces:  true,
		printTokens:     false,
		numberOfTokens:  4,
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

type ProxyPrinter struct {
	client          *scryfall.Client
	lang            scryfall.Lang
	deck            Deck
	printFrontFaces bool
	printBackFaces  bool
	printTokens     bool
	numberOfTokens  int
}

func (p *ProxyPrinter) WriteImageProxiesToFile(fileStr string) error {
	pdfFile, err := os.Create(fileStr)
	if err != nil {
		return err
	}
	err = p.WriteImageProxies(pdfFile)
	if err != nil {
		return err
	}
	err = pdfFile.Close()
	if err != nil {
		return err
	}
	return nil
}

func (p *ProxyPrinter) WriteImageProxies(w io.Writer) error {
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

	deck := p.collectProxyDeck()
	for _, s := range deck.Sections {
		if p.printSection(s.Name) {
			writeSection(s.Cards)
		}
	}
	return pdf.Output(w)
}

func (p *ProxyPrinter) WriteTextProxiesToFile(fileStr string) error {
	pdfFile, err := os.Create(fileStr)
	if err != nil {
		return err
	}
	err = p.WriteTextProxies(pdfFile)
	if err != nil {
		return err
	}
	err = pdfFile.Close()
	if err != nil {
		return err
	}
	return nil
}

func (p *ProxyPrinter) WriteTextProxies(w io.Writer) error {
	pdf := gofpdf.New("L", "mm", "A4", "")

	pdf.SetFont("Arial", "", 8)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFillColor(255, 255, 255)

	rep := strings.NewReplacer("−", "-", "\n", "\n\n")
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	writeSection := func(cards []Card) {
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
			pdf.SetFont("Arial", "B", 10)
			pdf.CellFormat(cardWidth-2*2, 6, tr(card.Name), "", 0, "LM", false, 0, "")
			pdf.SetFont("Arial", "", 8)

			if card.ManaCost != "" {
				pdf.MoveTo(x+2, y+2)
				pdf.CellFormat(cardWidth-2*2, 6, tr(card.ManaCost), "", 0, "RM", false, 0, "")
			}

			pdf.Line(x+2, y+2+6, x+cardWidth-2, y+2+6)

			pdf.MoveTo(x+2, y+2+6)
			pdf.CellFormat(cardWidth-2*2, 6, tr(card.TypeLine), "", 0, "LM", false, 0, "")

			pdf.MoveTo(x+2, y+2+6+6+1)
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

	deck := p.collectProxyDeck()
	for _, s := range deck.Sections {
		if p.printSection(s.Name) {
			writeSection(s.Cards)
		}
	}

	return pdf.Output(w)
}

func (p *ProxyPrinter) collectProxyDeck() Deck {
	versionFromCard := func(sc *scryfall.Card) *Version {
		if sc.Set != "" && sc.CollectorNumber != "" {
			return &Version{Set: sc.Set, CollectorNumber: sc.CollectorNumber}
		}
		return nil
	}

	cardFromCard := func(sc *scryfall.Card) Card {
		c := Card{
			Name:       sc.Name,
			ManaCost:   sc.ManaCost,
			TypeLine:   sc.TypeLine,
			OracleText: sc.OracleText,
			Power:      sc.Power,
			Toughness:  sc.Toughness,
			Loyalty:    sc.Loyalty,
			Version:    versionFromCard(sc),
		}
		if data, err := p.client.ImageByURL(sc.ImageURIs["large"]); err == nil {
			c.ImageData = data
		}
		return c
	}

	cardFromFace := func(sc *scryfall.CardFace) Card {
		c := Card{
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

	d := Deck{}
	frontFaces := Section{Name: FrontFaces}
	backFaces := Section{Name: BackFaces}
	tokens := Section{Name: Tokens}

	cards := p.deck.Cards()
	for _, card := range cards {
		sc := p.client.CardByName(card.Name)
		if sc == nil {
			continue
		}
		if p.lang != scryfall.LangEnglish {
			sc = p.client.CardBySetAndNumber(sc.Set, sc.CollectorNumber, p.lang)
		}
		if sc == nil {
			continue
		}

		switch sc.Layout {
		case scryfall.LayoutTransform:
			ff := cardFromFace(sc.Front())
			ff.Version = versionFromCard(sc)
			frontFaces.Cards = append(frontFaces.Cards, ff)
			bf := cardFromFace(sc.Back())
			bf.Version = versionFromCard(sc)
			backFaces.Cards = append(backFaces.Cards, bf)
		default:
			frontFaces.Cards = append(frontFaces.Cards, cardFromCard(sc))
		}

		if len(sc.AllParts) > 0 {
			for _, part := range sc.AllParts {
				if part.Component == "token" {
					if tc := p.client.CardByURL(part.URI); tc != nil {
						t := cardFromCard(tc)
						if !tokens.Cards.Contains(CardByName(t.Name)) {
							for i := 0; i < p.numberOfTokens; i++ {
								tokens.Cards = append(tokens.Cards, t)
							}
						}
					}
				}
			}
		}
	}
	if len(frontFaces.Cards) > 0 {
		d.Sections = append(d.Sections, frontFaces)
	}
	if len(backFaces.Cards) > 0 {
		d.Sections = append(d.Sections, backFaces)
	}
	if len(tokens.Cards) > 0 {
		d.Sections = append(d.Sections, tokens)
	}
	return d
}

func (p *ProxyPrinter) printSection(name string) bool {
	switch name {
	case FrontFaces:
		return p.printFrontFaces
	case BackFaces:
		return p.printBackFaces
	case Tokens:
		return p.printTokens
	}
	return true
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

const (
	FrontFaces = "FrontFaces"
	BackFaces  = "BackFaces"
	Tokens     = "Tokens"
)
