package mtg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/cognicraft/archive"
)

/* https://scryfall.com/docs/api/cards/named */

func NewScryfall(archive *archive.Archive) *Scryfall {
	return &Scryfall{
		baseURL: "https://api.scryfall.com",
		archive: archive,
	}
}

type Scryfall struct {
	baseURL string
	archive *archive.Archive
}

func (s *Scryfall) LargeImage(name string) ([]byte, error) {
	key := "/images/" + name
	imgRes, err := s.archive.Load(key)
	if err == nil {
		return imgRes.Data, nil
	}
	img, err := http.Get(s.urlLargeImageByName(name))
	if err != nil {
		return nil, err
	}
	defer img.Body.Close()
	bs, err := ioutil.ReadAll(img.Body)
	if err != nil {
		return nil, err
	}
	s.archive.Store(archive.JPEG(key, bs))
	return bs, err
}

func (s *Scryfall) Card(name string) (ScryfallCard, error) {
	resp, err := http.Get(s.urlCardByName(name))
	if err != nil {
		return ScryfallCard{}, err
	}
	defer resp.Body.Close()

	sc := ScryfallCard{}
	err = json.NewDecoder(resp.Body).Decode(&sc)
	if err != nil {
		return ScryfallCard{}, err
	}
	return sc, nil
}

func (s *Scryfall) urlLargeImageByName(name string) string {
	return fmt.Sprintf("%s/cards/named?format=image&version=large&fuzzy=%s", s.baseURL, url.QueryEscape(name))
}

func (s *Scryfall) urlCardByName(name string) string {
	return fmt.Sprintf("%s/cards/named?fuzzy=%s", s.baseURL, url.QueryEscape(name))
}

type ScryfallCard struct {
	Object          string            `json:"object"`
	ID              string            `json:"id"`
	OracleID        string            `json:"oracle_id"`
	MultiverseIDs   []int             `json:"multiverse_ids"`
	MtgoID          int               `json:"mtgo_id"`
	MtgoFoilID      int               `json:"mtgo_foil_id"`
	TcgPlayerID     int               `json:"tcgplayer_id"`
	Name            string            `json:"name"`
	Lang            string            `json:"lang"`
	ReleasedAt      string            `json:"released_at"`
	URI             string            `json:"uri"`
	ScryfallURI     string            `json:"scryfall_uri"`
	Layout          string            `json:"layout"`
	HighresImage    bool              `json:"highres_image"`
	ImageURIs       map[string]string `json:"image_uris,omitempty"`
	ManaCost        string            `json:"mana_cost"`
	CMC             float64           `json:"cmc"`
	TypeLine        string            `json:"type_line"`
	OracleText      string            `json:"oracle_text"`
	Colors          []string          `json:"colors"`
	ColorIdentity   []string          `json:"color_identity"`
	Legalities      map[string]string `json:"legalities"`
	Games           []string          `json:"games"`
	Reserved        bool              `json:"reserved"`
	Foil            bool              `json:"foil"`
	Nonfoil         bool              `json:"nonfoil"`
	Oversized       bool              `json:"oversized"`
	Promo           bool              `json:"promo"`
	Reprint         bool              `json:"reprint"`
	Variation       bool              `json:"variation"`
	Set             string            `json:"set"`
	SetName         string            `json:"set_name"`
	SetType         string            `json:"set_tye"`
	SetURI          string            `json:"set_uri"`
	SetSearchURI    string            `json:"set_search_uri"`
	ScryfallSetURI  string            `json:"scryfall_set_uri"`
	RulingsURI      string            `json:"rulings_uri"`
	PrintsSearchURI string            `json:"prints_search_uri"`
	CollectorNumber string            `json:"collector_number"`
	Digital         bool              `json:"digital"`
	Rarity          string            `json:"rarity"`
	CardBackID      string            `json:"card_back_id"`
	Artist          string            `json:"artist"`
	ArtistIDs       []string          `json:"artist_ids"`
	IllustrationID  string            `json:"illustration_id"`
	BorderColor     string            `json:"border_color"`
	Frame           string            `json:"frame"`
	FullArt         bool              `json:"full_art"`
	Textless        bool              `json:"textless"`
	Booster         bool              `json:"booster"`
	StorySpotlight  bool              `json:"story_spotlight"`
	EDHRecRank      int               `json:"edhrec_rank"`
	Prices          map[string]string `json:"prices"`
	RelatedURIs     map[string]string `json:"related_uris"`
	PurchaseURIs    map[string]string `json:"purchase_uris"`
}
