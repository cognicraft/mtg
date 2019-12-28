package scryfall

import (
	"fmt"

	"github.com/cognicraft/archive"
)

// Card objects represent individual Magic: The Gathering cards that players
// could obtain and add to their collection (with a few minor exceptions).
type Card struct {
	client *Client

	// This card’s Arena ID, if any.
	// A large percentage of cards are not available on Arena and do not have this ID.
	ArenaID int `json:"arena_id,omitempty"`

	// A unique ID for this card in Scryfall’s database.
	ID string `json:"id"`

	// A language code for this printing.
	Lang Lang `json:"lang"`

	// This card’s Magic Online ID (also known as the Catalog ID), if any.
	// A large percentage of cards are not available on Magic Online and do not have this ID.
	MtgoID int `json:"mtgo_id,omitempty"`

	// This card’s foil Magic Online ID (also known as the Catalog ID), if any.
	// A large percentage of cards are not available on Magic Online and do not have this ID.
	MtgoFoilID int `json:"mtgo_foil_id,omitempty"`

	// This card’s multiverse IDs on Gatherer, if any, as an array of integers.
	// Note that Scryfall includes many promo cards, tokens, and other esoteric objects that
	// do not have these identifiers.
	MultiverseIDs []int `json:"multiverse_ids,omitempty"`

	// This card’s ID on TCGplayer’s API, also known as the productId.
	TcgPlayerID int `json:"tcgplayer_id"`

	// A content type for this object, always card.
	Object string `json:"object"`

	// A unique ID for this card’s oracle identity. This value is consistent across reprinted card editions,
	// and unique among different cards with the same name (tokens, Unstable variants, etc).
	OracleID string `json:"oracle_id"`

	// A link to where you can begin paginating all re/prints for this card on Scryfall’s API.
	PrintsSearchURI string `json:"prints_search_uri"`

	// A link to this card’s rulings list on Scryfall’s API.
	RulingsURI string `json:"rulings_uri"`

	// A link to this card’s permapage on Scryfall’s website.
	ScryfallURI string `json:"scryfall_uri"`

	// A link to this card object on Scryfall’s API.
	URI string `json:"uri"`

	// If this card is closely related to other cards, this property will be an array with Related Card Objects.
	AllParts []*RelatedCard `json:"all_parts,omitempty"`

	// An array of Card Face objects, if this card is multifaced.
	CardFaces []*CardFace `json:"card_faces,omitempty"`

	// The card’s converted mana cost. Note that some funny cards have fractional mana costs.
	CMC float64 `json:"cmc"`

	// This card’s colors, if the overall card has colors defined by the rules.
	// Otherwise the colors will be on the card_faces objects, see below.
	Colors []string `json:"colors,omitempty"`

	// This card’s color identity.
	ColorIdentity []string `json:"color_identity"`

	// The colors in this card’s color indicator, if any.
	// A nil value for this field indicates the card does not have one.
	ColorIndicator []string `json:"color_indicator,omitempty"`

	// This card’s overall rank/popularity on EDHREC. Not all cards are ranked.
	EDHRecRank int `json:"edhrec_rank,omitempty"`

	// True if this printing exists in a foil version.
	Foil bool `json:"foil"`

	// This card’s hand modifier, if it is Vanguard card. This value will contain a delta, such as -1.
	HandModifier string `json:"hand_modifier,omitempty"`

	// 	A code for this card’s layout.
	Layout Layout `json:"layout"`

	// An object describing the legality of this card across play formats.
	// Possible legalities are legal, not_legal, restricted, and banned.
	Legalities map[string]Legality `json:"legalities"`

	// This card’s life modifier, if it is Vanguard card. This value will contain a delta, such as +2.
	LifeModifier string `json:"life_modifier,omitempty"`

	// This loyalty if any. Note that some cards have loyalties that are not numeric, such as X.
	Loyalty string `json:"loyalty,omitempty"`

	// The mana cost for this card. This value will be any empty string "" if the cost is absent.
	// Remember that per the game rules, a missing mana cost and a mana cost of {0} are different values.
	// Multi-faced cards will report this value in card faces.
	ManaCost string `json:"mana_cost,omitemtpy"`

	// The name of this card. If this card has multiple faces, this field will contain both names separated by ␣//␣.
	Name string `json:"name"`

	// True if this printing exists in a nonfoil version.
	Nonfoil bool `json:"nonfoil"`

	// The Oracle text for this card, if any.
	OracleText string `json:"oracle_text,omitempty"`

	// True if this card is oversized.
	Oversized bool `json:"oversized"`

	// This card’s power, if any. Note that some cards have powers that are not numeric, such as *.
	Power string `json:"power,omitempty"`

	// True if this card is on the Reserved List.
	Reserved bool `json:"reserved"`

	// This card’s toughness, if any. Note that some cards have toughnesses that are not numeric, such as *.
	Toughness string `json:"toughness,omitempty"`

	// The type line of this card.
	TypeLine string `json:"type_line"`

	// The name of the illustrator of this card. Newly spoiled cards may not have this field yet.
	Artist string `json:"artist,omitempty"`

	// Whether this card is found in boosters.
	Booster bool `json:"booster"`

	// This card’s border color: black, borderless, gold, silver, or white.
	BorderColor string `json:"border_color"`

	// The Scryfall ID for the card back design present on this card.
	CardBackID string `json:"card_back_id"`

	// This card’s collector number. Note that collector numbers can contain non-numeric characters, such as letters or ★.
	CollectorNumber string `json:"collector_number"`

	// True if this is a digital card on Magic Online.
	Digital bool `json:"digital"`

	// The flavor text, if any.
	FlavorText string `json:"flavor_text,omitempty"`

	// This card’s frame effects, if any.
	FrameEffects []string `json:"frame_effects,omitempty"`

	// This card’s frame layout.
	Frame Frame `json:"frame"`

	// True if this card’s artwork is larger than normal.
	FullArt bool `json:"full_art"`

	// A list of games that this card print is available in, paper, arena, and/or mtgo.
	Games []string `json:"games"`

	// True if this card’s imagery is high resolution.
	HighresImage bool `json:"highres_image"`

	// A unique identifier for the card artwork that remains consistent across reprints.
	// Newly spoiled cards may not have this field yet.
	IllustrationID string `json:"illustration_id,omitempty"`

	// An object listing available imagery for this card. See the Card Imagery article for more information.
	ImageURIs map[string]string `json:"image_uris,omitempty"`

	// An object containing daily price information for this card, including usd, usd_foil,
	// eur, and tix prices, as strings.
	Prices map[string]string `json:"prices"`

	// The localized name printed on this card, if any.
	PrintedName string `json:"printed_name,omitempty"`

	// The localized text printed on this card, if any.
	PrintedText string `json:"printed_text,omitempty"`

	// The localized type line printed on this card, if any.
	PrintedTypeLine string `json:"printed_type_line,omitempty"`

	// True if this card is a promotional print.
	Promo bool `json:"promo"`

	// An array of strings describing what categories of promo cards this card falls into.
	PromoTypes []string `json:"promo_types,omitempty"`

	// An object providing URIs to this card’s listing on major marketplaces.
	PurchaseURIs map[string]string `json:"purchase_uris"`

	// This card’s rarity. One of common, uncommon, rare, or mythic.
	Rarity string `json:"rarity"`

	// An object providing URIs to this card’s listing on other Magic: The Gathering online resources.
	RelatedURIs map[string]string `json:"related_uris"`

	// The date this card was first released.
	ReleasedAt string `json:"released_at"`

	// True if this card is a reprint.
	Reprint bool `json:"reprint"`

	// A link to this card’s set on Scryfall’s website.
	ScryfallSetURI string `json:"scryfall_set_uri"`

	// This card’s full set name.
	SetName string `json:"set_name"`

	// A link to where you can begin paginating this card’s set on the Scryfall API.
	SetSearchURI string `json:"set_search_uri"`

	// The type of set this printing is in.
	SetType string `json:"set_type"`

	// A link to this card’s set object on Scryfall’s API.
	SetURI string `json:"set_uri"`

	// This card’s set code.
	Set string `json:"set"`

	// True if this card is a Story Spotlight.
	StorySpotlight bool `json:"story_spotlight"`

	// True if the card is printed without text.
	Textless bool `json:"textless"`

	// Whether this card is a variation of another printing.
	Variation bool `json:"variation"`

	// The printing ID of the printing this card is a variation of.
	VariationOf string `json:"variation_of,omitempty"`

	// This card’s watermark, if any.
	Watermark string `json:"watermark,omitempty"`

	Preview *Preview `json:"preview,omitempty"`
}

func (c *Card) Printings() *List {
	l := List{}
	err := c.client.doGetJSON(c.PrintsSearchURI, &l)
	if err != nil {
		c.client.logf("%v", err)
		return nil
	}
	l.client = c.client
	return &l
}

func (c *Card) Image(version string) ([]byte, error) {
	c.client.logf("[DEBUG] Card.Image(%q)", version)
	url, ok := c.ImageURIs[version]
	if !ok {
		return nil, fmt.Errorf("unknown version")
	}
	img, err := c.client.cache.Load(url)
	if err == nil {
		c.client.logf("[DEBUG]   retrieved from cache")
		return img.Data, nil
	}
	data, err := c.client.doGetBytes(url)
	if err != nil {
		c.client.logf("[ERROR]   %v", err)
		return nil, err
	}
	c.client.cache.Store(archive.JPEG(url, data))
	c.client.logf("[DEBUG]   retrieved from scryfall")
	return data, nil
}

func (c *Card) Front() *CardFace {
	if len(c.CardFaces) > 0 {
		cf := c.CardFaces[0]
		cf.client = c.client
		return cf
	}
	return nil
}

func (c *Card) Back() *CardFace {
	if len(c.CardFaces) > 1 {
		cf := c.CardFaces[1]
		cf.client = c.client
		return cf
	}
	return nil
}

func (c *Card) IsLegalIn(format string) bool {
	if c == nil {
		return false
	}
	return c.Legalities[format] == LegalityLegal
}

type CardFace struct {
	client *Client

	// The name of the illustrator of this card face. Newly spoiled cards may not have this field yet.
	Artist string `json:"artist,omitempty"`

	// The colors in this face’s color indicator, if any.
	ColorIndicator []Color `json:"color_indicator,omitempty"`

	// Colors is this face’s colors.
	Colors []Color `json:"colors,omitempty"`

	// FlavorText is the flavor text printed on this face, if any.
	FlavorText string `json:"flavor_text,omitempty"`

	// IllustrationID is a unique identifier for the card face artwork that
	// remains consistent across reprints. Newly spoiled cards may not have
	// this field yet.
	IllustrationID string `json:"illustration_id,omitempty"`

	// ImageURIs is an object providing URIs to imagery for this face, if
	// this is a double-sided card. If this card is not double-sided, then the
	// image_uris property will be part of the parent object instead.
	ImageURIs map[string]string `json:"image_uris,omitempty"`

	// Loyalty is this face’s loyalty, if any.
	Loyalty string `json:"loyalty,omitempty"`

	// ManaCost is the mana cost for this face. This value will be any
	// empty string "" if the cost is absent. Remember that per the game
	// rules, a missing mana cost and a mana cost of {0} are different values.
	ManaCost string `json:"mana_cost"`

	// Name is the name of this particular face.
	Name string `json:"name"`

	// A content type for this object, always card_face.
	Object string `json:"object"`

	// OracleText is the Oracle text for this face, if any.
	OracleText string `json:"oracle_text,omitempty"`

	// Power is this face’s power, if any. Note that some cards have powers
	// that are not numeric, such as *.
	Power string `json:"power,omitempty"`

	// PrintedName is the printed name of this particular face.
	// This will only be set if the card is not in English.
	PrintedName string `json:"printed_name,omitempty"`

	// PrintedText is the printed text for this face, if any.
	// This will only be set if the card is not in English.
	PrintedText string `json:"printed_text,omitempty"`

	//The localized type line printed on this face, if any.
	PrintedTypeLine string `json:"printed_type_line,omitempty"`

	// Toughness is this face’s toughness, if any.
	Toughness string `json:"toughness,omitempty"`

	// TypeLine is the type line of this particular face.
	TypeLine string `json:"type_line"`

	// The watermark on this particulary card face, if any.
	Watermark string `json:"watermark,omitempty"`
}

func (cf *CardFace) Image(version string) ([]byte, error) {
	cf.client.logf("[DEBUG] CardFace.Image(%q)", version)
	url, ok := cf.ImageURIs[version]
	if !ok {
		return nil, fmt.Errorf("unknown version")
	}
	img, err := cf.client.cache.Load(url)
	if err == nil {
		cf.client.logf("[DEBUG]   retrieved from cache")
		return img.Data, nil
	}
	data, err := cf.client.doGetBytes(url)
	if err != nil {
		cf.client.logf("[ERROR]   %v", err)
		return nil, err
	}
	cf.client.cache.Store(archive.JPEG(url, data))
	cf.client.logf("[DEBUG]   retrieved from scryfall")
	return data, nil
}

// Cards that are closely related to other cards (because they call them by name, or generate a token, or meld, etc)
// have a all_parts property that contains Related Card objects.
type RelatedCard struct {
	client *Client

	// An unique ID for this card in Scryfall’s database.
	ID string `json:"id"`

	// A content type for this object, always related_card.
	Object string `json:"object"`

	// A field explaining what role this card plays in this relationship,
	// one of token, meld_part, meld_result, or combo_piece.
	Component string `json:"component"`

	// The name of this particular related card.
	Name string `json:"name"`

	// The type line of this card.
	TypeLine string `json:"type_line"`

	// A URI where you can retrieve a full object describing this card on Scryfall’s API.
	URI string `json:"uri"`
}

type Lang string

const (
	LangEnglish            Lang = "en"
	LangSpanish            Lang = "es"
	LangFrench             Lang = "fr"
	LangGerman             Lang = "de"
	LangItalian            Lang = "it"
	LangPortuguese         Lang = "pt"
	LangJapanese           Lang = "ja"
	LangKorean             Lang = "ko"
	LangRussian            Lang = "ru"
	LangSimplifiedChinese  Lang = "zhs"
	LangTraditionalChinese Lang = "zht"
	LangHebrew             Lang = "he"
	LangLatin              Lang = "la"
	LangAncientGreek       Lang = "grc"
	LangArabic             Lang = "ar"
	LangSanskrit           Lang = "sa"
	LangPhyrexian          Lang = "px"
)

type Layout string

const (
	LayoutNormal           Layout = "normal"
	LayoutSplit            Layout = "split"
	LayoutFlip             Layout = "flip"
	LayoutTransform        Layout = "transform"
	LayoutMeld             Layout = "meld"
	LayoutLeveler          Layout = "leveler"
	LayoutSaga             Layout = "saga"
	LayoutPlanar           Layout = "planar"
	LayoutScheme           Layout = "scheme"
	LayoutVanguard         Layout = "vanguard"
	LayoutToken            Layout = "token"
	LayoutDoubleFacedToken Layout = "double_faced_token"
	LayoutEmblem           Layout = "emblem"
	LayoutAugment          Layout = "augment"
	LayoutHost             Layout = "host"
)

type Legality string

const (
	LegalityLegal      Legality = "legal"
	LegalityNotLegal   Legality = "not_legal"
	LegalityBanned     Legality = "banned"
	LegalityRestricted Legality = "restricted"
)

type Frame string

const (
	Frame1993   Frame = "1993"
	Frame1997   Frame = "1997"
	Frame2003   Frame = "2003"
	Frame2015   Frame = "2015"
	FrameFuture Frame = "future"
)

type FrameEffect string

const (
	FrameEffectLegendary      FrameEffect = "legendary"
	FrameEffectMiracle        FrameEffect = "miracle"
	FrameEffectNyxTouched     FrameEffect = "nyxtouched"
	FrameEffectDraft          FrameEffect = "draft"
	FrameEffectDevoid         FrameEffect = "devoid"
	FrameEffectTombstone      FrameEffect = "tombstone"
	FrameEffectColorShifted   FrameEffect = "colorshifted"
	FrameEffectSunMoonDFC     FrameEffect = "sunmoondfc"
	FrameEffectCompassLandDFC FrameEffect = "compasslanddfc"
	FrameEffectOriginPWDFC    FrameEffect = "originpwdfc"
	FrameEffectMoonEldraziDFC FrameEffect = "mooneldrazidfc"
)

type Preview struct {

	// The date this card was previewed.
	PreviewedAt string `json:"previewed_at"`

	// A link to the preview for this card.
	SourceURI string `json:"source_uri"`

	// The name of the source that previewed this card.
	Source string `json:"source"`
}

type Component string

const (
	ComponentToken      Component = "token"
	ComponentMeldPart   Component = "meld_part"
	ComponentMeldResult Component = "meld_result"
	ComponentComboPiece Component = "combo_piece"
)

type Color string
