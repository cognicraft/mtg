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
	ImageURIs map[string]string `json:"image_uris,omitempty"`
}
