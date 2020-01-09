package scryfall

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/cognicraft/archive"
)

/* https://scryfall.com/docs/api/cards/named */
/* /cards/:code/:number(/:lang)*/

func New(opts ...func(*Client) error) (*Client, error) {
	c := &Client{
		baseURL:    "https://api.scryfall.com",
		lang:       LangEnglish,
		delay:      100 * time.Millisecond,
		httpClient: &http.Client{},
	}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	if c.logf == nil {
		c.logf = discardf
	}
	if c.cache == nil {
		arc, err := archive.Open(":memory:")
		if err != nil {
			return nil, err
		}
		c.cache = arc
	}
	return c, nil
}

func Cache(cache *archive.Archive) func(*Client) error {
	return func(s *Client) error {
		s.cache = cache
		return nil
	}
}

func Language(l Lang) func(*Client) error {
	return func(s *Client) error {
		s.lang = l
		return nil
	}
}

func Debug(c *Client) error {
	c.logf = func(format string, args ...interface{}) {
		fmt.Printf(format+"\n", args...)
	}
	return nil
}

type Client struct {
	baseURL    string
	cache      *archive.Archive
	lang       Lang
	logf       func(string, ...interface{})
	delay      time.Duration
	lastAccess time.Time
	httpClient *http.Client
}

func (c *Client) CardByName(name string) *Card {
	c.logf("[DEBUG] Card(%q)", name)

	url := c.urlCardByName(name)

	card := Card{}
	if err := archive.LoadJSON(c.cache, url, &card); err == nil {
		c.logf("[DEBUG]   retrieved from cache")
		return &card
	}

	err := c.doGetJSON(url, &card)
	if err != nil {
		c.logf("[ERROR]   %v", err)
		return nil
	}
	c.cache.Store(archive.GenericJSON(url, card))
	c.logf("[DEBUG]   retrieved from scryfall")
	return &card
}

func (c *Client) CardByURL(url string) *Card {
	card := Card{}
	if err := archive.LoadJSON(c.cache, url, &card); err == nil {
		c.logf("[DEBUG]   retrieved from cache")
		return &card
	}

	err := c.doGetJSON(url, &card)
	if err != nil {
		c.logf("[ERROR]   %v", err)
		return nil
	}
	c.cache.Store(archive.GenericJSON(url, card))
	c.logf("[DEBUG]   retrieved from scryfall")
	return &card
}

func (s *Client) ImageByURL(url string) ([]byte, error) {
	s.logf("[DEBUG] Image(%q)", url)
	img, err := s.cache.Load(url)
	if err == nil {
		s.logf("[DEBUG]   retrieved from cache")
		return img.Data, nil
	}
	data, err := s.doGetBytes(url)
	if err != nil {
		s.logf("[ERROR]   %v", err)
		return nil, err
	}
	s.cache.Store(archive.JPEG(url, data))
	s.logf("[DEBUG]   retrieved from scryfall")
	return data, nil
}

func (c *Client) doGetJSON(url string, v interface{}) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	res, err := c.doRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s - %s", url, res.Status)
	}
	return json.NewDecoder(res.Body).Decode(v)
}

func (c *Client) doGetBytes(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	res, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s - %s", url, res.Status)
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) doRequest(req *http.Request) (*http.Response, error) {
	dt := time.Since(c.lastAccess)
	if dt < c.delay {
		time.Sleep(c.delay - dt)
	}
	c.lastAccess = time.Now()
	return c.httpClient.Do(req)
}

func (s *Client) urlCardByName(name string) string {
	return fmt.Sprintf("%s/cards/named?fuzzy=%s", s.baseURL, url.QueryEscape(name))
}

func (s *Client) urlCardImageByName(name string) string {
	return fmt.Sprintf("%s/cards/named?format=image&version=large&fuzzy=%s", s.baseURL, url.QueryEscape(name))
}

func discardf(string, ...interface{}) {}
