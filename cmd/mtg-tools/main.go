package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/cognicraft/archive"
	"github.com/cognicraft/hyper"
	"github.com/cognicraft/mtg"
	"github.com/cognicraft/mtg/scryfall"
	"github.com/cognicraft/mux"
)

func main() {

	bindFlag := flag.String("bind", ":8888", "Bind")
	cacheFlag := flag.String("cache", "cache.arc", "Cache")
	flag.Parse()

	cache, err := archive.Open(*cacheFlag)
	if err != nil {
		log.Fatal(err)
	}
	defer cache.Close()

	var scOpts []scryfall.ClientOption
	scOpts = append(scOpts, scryfall.Cache(cache))

	scry, err := scryfall.New(scOpts...)
	if err != nil {
		log.Fatal(err)
	}

	service := &Service{
		Scryfall: scry,
	}

	chain := mux.NewChain()
	router := mux.New()
	router.Route("/css/style.css").GET(chain.ThenFunc(service.handleGETStyleCSS))
	router.Route("/").GET(chain.ThenFunc(service.handleGET))
	router.Route("/").POST(chain.ThenFunc(service.handlePOST))

	s := &http.Server{
		Addr:           *bindFlag,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("binding to: %s", s.Addr)
	log.Printf("router configuration:\n\n%s\n", mux.Tree(router.Route("/")))
	log.Fatal(s.ListenAndServe())
}

type Service struct {
	Scryfall *scryfall.Client
}

func (s *Service) handleGETStyleCSS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(hyper.HeaderContentType, "text/css")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, css)
}

func (s *Service) handleGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(hyper.HeaderContentType, "text/html")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, index)
}

func (s *Service) handlePOST(w http.ResponseWriter, r *http.Request) {
	cmd := hyper.ExtractCommand(r)
	switch cmd.Action {
	case "generate-proxies":
		name := cmd.Arguments.String("name")
		deckText := cmd.Arguments.String("deck")
		deck, err := mtg.ParseDeck(strings.NewReader(deckText))
		if err != nil {
			hyper.WriteError(w, http.StatusBadRequest, err)
			return
		}
		deck.Name = name

		w.Header().Set(hyper.HeaderContentType, "application/pdf")
		err = mtg.NewProxyPrinter(s.Scryfall, deck).WriteImageProxies(w)
		if err != nil {
			log.Printf("%#v", err)
		}
	}
}

const css = `
* {
	margin: 0;
	padding: 0;
	box-sizing: border-box;
}

body {
	font-family: "Roboto", Helvetica, Arial, sans-serif;
	font-weight: 100;
	font-size: 12px;
	color: #777;
	background: #4CAF50;
	margin: 4em;
}

.card {
	background: #F9F9F9;
	width: 75%;
	padding: 25px;
	margin: auto;
	margin-top: 3em;
	margin-bottom: 3em;
	border-radius: 10px;
	box-shadow: 0 0 20px 0 rgba(0, 0, 0, 0.2), 0 5px 5px 0 rgba(0, 0, 0, 0.24);
}
.card h1 {
	margin-bottom: 1em;
}

.card ul {
	margin-top: .75em;
	margin-left: 1em;
	margin-right:1em;
	margin-bottom: .75em;
}

.card li {
	margin-top: .5em;
	margin-bottom: .5em;
}

label {
	margin-right: .5em;
}
textarea {
	margin: auto;
	width: 99%;
	resize: none;
}
`

const index = `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<title>MTG - Proxy Deck Generator</title>
	<link rel="stylesheet" href="/css/style.css">
	<!--[if IE]>
		<script src="http://html5shiv.googlecode.com/svn/trunk/html5.js"></script>
	<![endif]-->
</head>

<body translate="no">
	<div class="card">
		<h1>MTG - Proxy Deck Generator</h1>
		<form action="/" method="POST">
			<input type="hidden" name="@action" value="generate-proxies">
			<label>Name:</label><input type="text" name="name"></input><br/>
			<label>Deck:</label><textarea name="deck" cols="80" rows="30"></textarea><br/>
			<input type="submit" value="Generate Proxies" />
		</form>
	</div>
	<div class="card">
		<h1>The Staples Binder Method</h1>
		<p>The Staples Binder Method can be used to to save some cash while playing multiple decks within a format. With this method you will need at max 4 original copies of any given card in your collection. To reduce the amount of effort this method should only be used for cards that have a value greater than a few dollars.</p>
		<ul>
			<li>Put all decks in the same sleeves (brand/color).</li>
			<li>Keep the sleeved staples in a binder.</li>
			<li>Proxy print the needed number of any given staple with the decks name as an overlay.</li>
			<li>Cut out the proxy.</li>
			<li>Sleve up the proxy together with a basic land.</li>
			<li>Keep the sleeved proxy with the rest of the cards in the deck.</li>
			<li>Before the game pull out all proxies and replace them with the original printings from the Stapels Binder.</li>
			<li>After the game replace all of the staples in your deck with the proxies that where pulled out before the game.</li>
			<li><em>Tip:</em> The name of the deck is important if you want to always know where the original printing currently is. If you find a proxy within your Staples Binder the name will tell you where your original is located. In case you're trying to replace a proxy from deck A and find a proxy of deck B in your Staples Binder you know that the original is currently located in deck B.</li>
			<li><em>Tip:</em> In a casual setting and if your playgroup does not mind, you could skip replacing the proxies to save some time and effort. In case you are double sleeving your cards you will not notice a difference in thickness between (proxy + basic land + outer-sleeve) and (inner-sleeve + staple + outer-sleeve).</li>
		</ul>
	</div>
</body>
</html>
`
