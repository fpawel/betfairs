package main

import (
	"github.com/go-chi/chi"
	"github.com/fpawel/betfairs/football"
	"log"
	"net/http"
	"os"
	"strings"

	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/fpawel/betfairs/aping"
	"github.com/fpawel/betfairs/aping/listMarketBook"
	"github.com/fpawel/betfairs/aping/listMarketCatalogue"
	"github.com/fpawel/betfairs/webclient"
	"io/ioutil"
	"strconv"
	"github.com/fpawel/betfairs/event2"
)

func daemon() {

	apingSession := aping.NewSession(adminBetfairUser, adminBetfairPass)
	fmt.Println(apingSession.GetSession())

	betfairReader := BetfairClient{
		Football:            new(football.GamesReader),
		ListMarketCatalogue: listMarketCatalogue.New(apingSession),
		ListMarketBook:      listMarketBook.New(apingSession),
	}

	router := chi.NewRouter()
	var websocketUpgrader = websocket.Upgrader{EnableCompression: true}

	fileServer(router, "/", http.Dir("assets"))


	router.Get("/football/games", func(w http.ResponseWriter, r *http.Request) {
		games, err := betfairReader.Football.Read()
		setJsonResult(w, games, err)
	})

	router.Get("/football/games2", func(w http.ResponseWriter, r *http.Request) {
		var tmp int32
		games, err := betfairReader.ReadFootballGames2(&tmp)
		setJsonResult(w, games, err)
	})

	router.Get("/football/games3", func(w http.ResponseWriter, r *http.Request) {
		var tmp int32
		games, err := betfairReader.ReadFootballGames3(&tmp)
		setJsonResult(w, games, err)
	})

	router.Get("/football/games4", func(w http.ResponseWriter, r *http.Request) {
		var tmp int32
		games, err := betfairReader.ReadFootballGames4(&tmp)
		setJsonResult(w, games, err)
	})


	router.Get("/football/live", func(w http.ResponseWriter, r *http.Request) {

		conn, err := websocketUpgrader.Upgrade(w, r, nil)
		check(err)
		runWebSocketFootballLive( conn, betfairReader.Football)
		conn.Close()
	})

	router.Get("/football", func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocketUpgrader.Upgrade(w, r, nil)
		check(err)
		runWebSocketFootball( conn, betfairReader)
		conn.Close()
	})

	router.Get("/redirect-betfair/*", func(w http.ResponseWriter, r *http.Request) {
		redirect(webclient.NewURL(chi.URLParam(r, "*") ), w, r)
	})


	router.Get("/event/{eventID}", func(w http.ResponseWriter, r *http.Request) {
		eventID, err := strconv.Atoi(chi.URLParam(r, "eventID"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		marketCatalogues, err := betfairReader.ListMarketCatalogue.Read(eventID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		home, away,ok := betfairReader.Football.TeamsByID(eventID)
		if !ok {
			http.Error(w, fmt.Sprintf("game not found: %d", eventID), http.StatusBadRequest)
			return
		}

		setJsonResult(w, event2.NewEvent(marketCatalogues, home, away), nil)
	})

	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", "8080")
	}
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}

// fileServer conveniently sets up a http.fileServer handler to serve
// static files from a http.FileSystem.
func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("fileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

func redirect(urlStr string, w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest(r.Method, urlStr, r.Body)
	if err != nil {
		log.Fatal(err)
	}
	for key, value := range r.Header {
		req.Header.Set(key, strings.Join(value, "; "))
	}

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	for key, value := range response.Header {
		w.Header().Set(key, strings.Join(value, "; "))
	}

	w.WriteHeader(response.StatusCode)
	body, err := ioutil.ReadAll(response.Body)
	if err == nil {
		_, err := w.Write(body[:])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func setCompressedJSON(w http.ResponseWriter, data interface{}) {
	gz, err := gzip.NewWriterLevel(w, gzip.DefaultCompression)
	check(err)
	defer gz.Close()
	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	encoder := json.NewEncoder(gz)
	encoder.SetIndent("", "    ")
	check(encoder.Encode(data))
}

func setJsonResult(w http.ResponseWriter, data interface{}, err error) {

	if err != nil {
		var y struct {
			Error string `json:"error"`
		}
		y.Error = err.Error()
		setCompressedJSON(w, &y)
		return
	}

	var y struct {
		Ok interface{} `json:"ok"`
	}
	y.Ok = data
	setCompressedJSON(w, &y)

}
