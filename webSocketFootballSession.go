package main

import (
	"github.com/gorilla/websocket"
	"fmt"
	"time"

	"heroku.com/betfairs/football2"
	"sync/atomic"
)


type webSocketFootballSession struct {
	conn            *websocket.Conn
	gamesReader *football2.GamesReader
}

func (x webSocketFootballSession) run() {
	x.conn.EnableWriteCompression(true)
	sendGames := make(chan football2.Games)
	var interruptReadGames int32
	go func() {

		for {
			games, err := x.gamesReader.Read(&interruptReadGames)
			if err == football2.ErrorInterrupted {
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
				x.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			x.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			changes := games.Changes(nextGames)

			if !changes.Empty() {
				err := x.conn.WriteJSON(changes)
				if err != nil {
					fmt.Println("WebSocket: error 1:", err)
					return
				}
			}
			games = nextGames
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
	atomic.AddInt32(&interruptReadGames, 1)
	close(sendGames)
	<- doneSendGames
}


