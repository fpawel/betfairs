package main

import (
	"net/http"
	"strings"
	"github.com/go-chi/chi"
	"log"
	"os"
	"heroku.com/betfairs/football"

	"io/ioutil"
	"heroku.com/betfairs/webclient"
	"encoding/json"
	"compress/gzip"
	"github.com/gorilla/websocket"
	"heroku.com/betfairs/aping/apingEvents"
	"heroku.com/betfairs/aping"
	"fmt"
)


func daemon (){

	apingSession := aping.NewSession(adminBetfairUser, adminBetfairPass)
	fmt.Println( apingSession.GetSession() )

	eventsReader := apingEvents.NewSyncReader(apingSession)
	footballReader := new(football.SyncReader)
	router := chi.NewRouter()
	var websocketUpgrader = websocket.Upgrader{EnableCompression: true}


	FileServer(router, "/", http.Dir("assets"))

	router.Get("/football/games", func(w http.ResponseWriter, r *http.Request) {
		games,err := footballReader.Read()
		setJsonResult(w, games, err)
	})

	router.Get("/football", func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocketUpgrader.Upgrade(w, r, nil)
		check(err)
		webSocketFootball(conn, footballReader, eventsReader)
		conn.Close()
	})

	router.Get("/redirect-betfair/*", func(w http.ResponseWriter, r *http.Request) {
		redirect(webclient.NewURL(chi.URLParam(r, "*")), w, r)
	})

	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", "8080")
	}
	log.Fatal(http.ListenAndServe(":" + os.Getenv("PORT"), router))
}



// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
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


func redirect (urlStr string , w http.ResponseWriter, r *http.Request) {
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

	if err != nil {
		log.Fatal(err)
	}
	defer gz.Close()

	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	encoder := json.NewEncoder(gz)
	encoder.SetIndent("", "    ")

	err = encoder.Encode(data)
	if err != nil {
		log.Fatal(err)
	}
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
