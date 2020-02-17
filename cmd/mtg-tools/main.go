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
		deckText := cmd.Arguments.String("deck")
		deck, err := mtg.ParseDeck(strings.NewReader(deckText))
		if err != nil {
			hyper.WriteError(w, http.StatusBadRequest, err)
			return
		}
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
	line-height: 30px;
	color: #777;
	background: #4CAF50;
	margin: 4em;
}

.card {
	background: #F9F9F9;
	width: 50%;
	padding: 25px;
	margin: auto auto;
	border-radius: 10px;
	box-shadow: 0 0 20px 0 rgba(0, 0, 0, 0.2), 0 5px 5px 0 rgba(0, 0, 0, 0.24);
}
.card h1 {
	margin-bottom: 1em;
}

textarea {
	margin: auto;
	width: 98%;
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
			<textarea name="deck" cols="80" rows="30"></textarea><br/>
			<input type="submit" value="Generate Proxies" />
		</form>
	</div>
</body>
</html>
`