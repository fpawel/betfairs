package main

import (
	"github.com/gorilla/websocket"
	"fmt"
	"time"

	"github.com/fpawel/betfairs/football/football2"
	"sync/atomic"
	"github.com/fpawel/betfairs/football/football4"
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
	atomic.AddInt32(&interruptReadGames, 1)
	close(sendGames)
	<- doneSendGames
}

func runWebSocketFootballPrices(conn *websocket.Conn, betfair BetfairClient) {
	conn.EnableWriteCompression(true)
	sendGames := make(chan []football4.Game)
	var interruptReadGames int32
	go func() {
		for {
			games4,err := betfair.ReadFootballGames4(&interruptReadGames)
			if err == ErrorInterrupted {
				return
			}
			if err != nil {
				fmt.Println("ERROR football:", err )
				continue
			}
			sendGames <- games4
		}
	}()

	doneSendGames := make(chan bool) // цикл записи завершён
	go func () {
		defer func() {
			doneSendGames <- true
		}()
		for {
			games4, ok := <-sendGames
			if !ok { // если канал send закрыт, прервать цикл записи
				conn.WriteMessage(websocket.CloseMessage,[]byte{})
				return
			}
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			err := conn.WriteJSON(games4)
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
	close(sendGames)
	<- doneSendGames
}
