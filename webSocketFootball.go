package main

import (
	"github.com/gorilla/websocket"
	"fmt"
	"time"

	"github.com/fpawel/betfairs/football/football2"
	"sync/atomic"
	"github.com/fpawel/betfairs/football"
	"github.com/fpawel/betfairs/aping"
)

func runWebSocketFootball(conn *websocket.Conn, betfair BetfairClient) {
	conn.EnableWriteCompression(true)
	sendGames := make(chan football2.Games)
	var interruptReadGames int32
	go func() {

		for {
			games, err := betfair.ReadFootballGames2(&interruptReadGames)
			if err == ErrorInterrupted {
				return
			}
			if err != nil {
				fmt.Println("ERROR football:", err )
				continue
			}
			sendGames <- games
		}
	}()

	doneSendGames := make(chan bool) // цикл записи завершён
	var games football2.Games
	go func () {
		defer func() {
			doneSendGames <- true
		}()
		for {
			nextGames, ok := <-sendGames
			if !ok { // если канал send закрыт, прервать цикл записи
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			changes := games.Changes(nextGames)

			if !changes.Empty() {
				err := conn.WriteJSON(changes)
				if err != nil {
					fmt.Println("WebSocket: error 1:", err)
					return
				}
			}
			games = nextGames
		}
	}()
	for {
		messageType, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				fmt.Println("WebSocket error 2:", err)
			}
			break
		}
		switch messageType {
		case websocket.CloseMessage:
			break
		}
	}
	atomic.AddInt32( &interruptReadGames, 1)
	close(sendGames)
	<- doneSendGames
}

func runWebSocketFootballLive(conn *websocket.Conn, rdr BetfairClient) {
	conn.EnableWriteCompression(true)
	done := make(chan bool)
	var interruptReadGames int32

	go func() {
		defer func() {
			done <- true
		}()
		for {
			var games []football.Game
			games, err := rdr.Football.Read()
			if atomic.LoadInt32(&interruptReadGames) > 0 {
				return
			}
			if err != nil {
				fmt.Println("ERROR football:", err )
				continue
			}

			var gamesLive []football.GameLive

			for _, game := range games {
				if !game.InPlay || !game.HasMinute() {
					continue
				}
				mc,err := rdr.ListMarketCatalogue.Read(game.ID)
				if err != nil {
					if err != aping.ErrorNoMarkets{
						fmt.Println("ERROR football 2:", err )
					}
					continue
				}
				gamesLive = append(gamesLive, game.GameLive(mc[0].Event.OpenDate))
			}

			err = conn.WriteJSON(gamesLive)
			if err != nil {
				fmt.Println("WebSocket: error 3:", err)
				return
			}
		}
	}()

	for {
		messageType, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				fmt.Println("WebSocket error 4:", err)
			}
			break
		}
		switch messageType {
		case websocket.CloseMessage:
			break
		}
	}
	atomic.AddInt32(&interruptReadGames, 1)
	<- done
}