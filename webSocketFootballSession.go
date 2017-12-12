package main

import (
	"github.com/gorilla/websocket"
	"fmt"
	"time"
	"sync/atomic"

	"heroku.com/betfairs/football2"
)


type webSocketFootballSession struct {
	conn            *websocket.Conn
	*football2.GamesReader
}

func (x webSocketFootballSession) run() {
	x.conn.EnableWriteCompression(true)
	x.conn.SetReadLimit(100000)
	x.conn.SetReadDeadline(time.Now().Add(20 * time.Second))
	x.conn.SetPongHandler(func(string) error {
		x.conn.SetReadDeadline(time.Now().Add(20 * time.Second))
		return nil
	})

	sendGames := make(chan football2.Games)
	//interruptReadGamesDelay := make(chan bool, 2) // прервать цикл поллинга футбола
	go func() {

		for {
			games, err := x.Read()
			if err != nil {
				fmt.Println("ERROR football:", err )
				continue
			}
			if atomic.LoadInt32(&x.DoneFlag) > 0 {
				return
			}
			sendGames <- games
			//select {
			//case <-interruptReadGamesDelay:
			//	return
			//case <-time.After(5 * time.Second):
			//	continue
			//}
		}
	}()



	pingTicker := time.NewTicker(5 * time.Second) // пинговать клиента раз в 5 секунд
	done := make(chan bool) // цикл записи завершён
	var games football2.Games
	go func () {
		defer func() {
			atomic.AddInt32(&x.DoneFlag, 1)
			done <- true
		}()
		for {
			select {
			case nextGames, ok := <-sendGames:
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

			case <-pingTicker.C:
				x.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
				err := x.conn.WriteMessage(websocket.PingMessage, []byte{})
				if err != nil {
					fmt.Println("WebSocket: error 2:", err)
					return
				}
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




	atomic.AddInt32(&x.DoneFlag, 1)
	pingTicker.Stop()
	close(sendGames)
	//interruptReadGamesDelay <- true
	<- done
}


