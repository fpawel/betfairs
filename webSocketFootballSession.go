package main

import (
	"github.com/gorilla/websocket"
	"fmt"
	"time"
	"heroku.com/betfairs/football"
	"sync/atomic"

	"heroku.com/betfairs/aping/listMarketCatalogue"
	"heroku.com/betfairs/countries"
	"heroku.com/betfairs/aping/listMarketBook"
)


type webSocketFootballSession struct {
	conn            *websocket.Conn
	football        *football.SyncReader
	marketCatalogue *listMarketCatalogue.Reader
	marketBook      *listMarketBook.Reader
}

func (x webSocketFootballSession) run() {
	x.conn.EnableWriteCompression(true)
	x.conn.SetReadLimit(100000)
	x.conn.SetReadDeadline(time.Now().Add(20 * time.Second))
	x.conn.SetPongHandler(func(string) error {
		x.conn.SetReadDeadline(time.Now().Add(20 * time.Second))
		return nil
	})

	type game struct {
		football.Game
		Competition string `json:"competition"`
		Country string `json:"country"`
		MainPrices [6]float64 `json:"main_prices,omitempty"`
	}

	pingTicker := time.NewTicker(5 * time.Second) // пинговать клиента раз в 5 секунд
	sendGames := make(chan []game)
	done := make(chan bool) // цикл записи завершён

	var doneFlag int32

	go func () {
		defer func() {
			atomic.AddInt32(&doneFlag, 1)
			done <- true
		}()
		for {
			select {
			case games, ok := <-sendGames:
				if !ok { // если канал send закрыт, прервать цикл записи
					x.conn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}

				x.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

				err := x.conn.WriteJSON(games)
				if err != nil {
					println("WebSocket: error 1:", err)
					return
				}

			case <-pingTicker.C:
				x.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
				err := x.conn.WriteMessage(websocket.PingMessage, []byte{})
				if err != nil {
					println("WebSocket: error 2:", err)
					return
				}
			}
		}
	}()


	interruptReadGamesDelay := make(chan bool, 2) // прервать цикл поллинга футбола

	go func() {

		for {
			xs, err := x.football.Read()
			if err != nil {
				fmt.Println("ERROR football:", err )
				continue
			}
			if atomic.LoadInt32(&doneFlag) > 0 {
				return
			}
			var games []game
			for _,g := range xs {
				var game game
				game.Game = g
				marketCatalogues,ok := x.marketCatalogue.Get(game.ID)
				if ok {
					game.Competition = marketCatalogues[0].Competition.Name
					c := countries.ByAlpha2(marketCatalogues[0].Event.CountryCode)
					if c != nil {
						game.Country = c.Name
					} else {
						game.Country = marketCatalogues[0].Event.CountryCode
					}
					mainMarket,ok := marketCatalogues.MainMarket()
					if ok {
						t := time.Second
						if !game.InPlay {
							t = time.Minute
						}
						mainMarketBook, err := x.marketBook.Read([]string{mainMarket.ID}, t)
						if ok {
							game.MainPrices = mainMarketBook[0].Prices6()
						} else {
							fmt.Println(err)
						}
					}
				}
				games = append(games, game)
			}

			sendGames <- games

			select {
			case <-interruptReadGamesDelay:
				return
			case <-time.After(5 * time.Second):
				continue
			}
		}
	}()

	for {
		messageType, _, err := x.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				fmt.Println("WebSocket error 1:", err)
			}
			break
		}
		switch messageType {
		case websocket.CloseMessage:
			break
		}
	}

	atomic.AddInt32(&doneFlag, 1)
	pingTicker.Stop()
	close(sendGames)
	interruptReadGamesDelay <- true
	<- done
}
