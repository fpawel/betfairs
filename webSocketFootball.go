package main

import (
	"github.com/gorilla/websocket"
	"fmt"
	"time"
	"heroku.com/betfairs/football"
	"sync/atomic"

	"heroku.com/betfairs/aping/apingEvents"
	"heroku.com/betfairs/countries"
)


func webSocketFootball(conn *websocket.Conn, footballReader *football.SyncReader, eventsReader *apingEvents.SyncReader) {
	conn.EnableWriteCompression(true)
	conn.SetReadLimit(100000)
	conn.SetReadDeadline(time.Now().Add(20 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(20 * time.Second))
		return nil
	})

	type game struct {
		football.Game
		Competition string `json:"competition"`
		Country string `json:"country"`
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
					conn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}

				conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

				err := conn.WriteJSON(games)
				if err != nil {
					println("WebSocket: error 1:", err)
					return
				}

			case <-pingTicker.C:
				conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
				err := conn.WriteMessage(websocket.PingMessage, []byte{})
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
			xs, err := footballReader.Read()
			if err != nil {
				fmt.Println("ERROR football:", err )
				continue
			}
			if atomic.LoadInt32(&doneFlag) > 0 {
				return
			}
			var games []game
			for _,x := range xs {
				var game game
				game.Game = x
				event := eventsReader.Event(x.ID)
				if event != nil {
					if event.Competition != nil {
						game.Competition = event.Competition.Name
					}
					if err == nil {
						c := countries.ByAlpha2(event.CountryCode)
						if c != nil {
							game.Country = c.Name
						}
					}
				}
				games = append(games, game)
			}

			sendGames <- games

			select {
			case <-interruptReadGamesDelay:
				return
			case <-time.After(10 * time.Second):
				continue
			}
		}
	}()

	for {
		messageType, _, err := conn.ReadMessage()
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