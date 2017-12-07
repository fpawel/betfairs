package main

import (
	"time"

	"github.com/gorilla/websocket"
)

type webSocketWorkerConfig struct {
	WriteWait, // Time allowed to write a message to the peer.
	PongWait, // Time allowed to read the next pong message from the peer.
	PingPeriod time.Duration // Send pings to peer with this period. Must be less than pongWait.
	MaxMessageSize int64 // Maximum message size allowed from peer.
}

type webSocketWorker struct {
	conn      *websocket.Conn
	send      <-chan []byte
	interrupt <-chan bool
	done      chan<- bool
	f         func([]byte)
	writeWait,
	pongWait time.Duration
}

func webSocketWorkerDefaultConfig() webSocketWorkerConfig {
	const PongWait = 60 * time.Second
	return webSocketWorkerConfig{
		PongWait:       PongWait,
		PingPeriod:     (PongWait * 9) / 10,
		WriteWait:      10 * time.Second,
		MaxMessageSize: 100000,
	}
}

func runWebSocketWorker(conn *websocket.Conn, send <-chan []byte, config webSocketWorkerConfig, f func([]byte)) (err error) {

	done := make(chan bool)
	interrupt := make(chan bool, 2)

	x := &webSocketWorker{
		conn:      conn,
		send:      send,
		interrupt: interrupt,
		done:      done,
		f:         f,
		writeWait: config.WriteWait,
		pongWait:  config.PongWait,
	}

	conn.SetReadLimit(config.MaxMessageSize)
	conn.SetReadDeadline(time.Now().Add(config.PongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(config.PongWait))
		return nil
	})

	go x.WritePump()
	err = x.ReadPump()
	interrupt <- true
	<-done
	close(interrupt)
	close(done)

	return
}

//var upgrader = websocket.Upgrader{EnableCompression: true}
//router.Get("/data", func(w http.ResponseWriter, r *http.Request) {
//	conn, err := upgrader.Upgrade(w, r, nil)
//	defer conn.Close()
//	check(err)
//	conn.EnableWriteCompression(true)
//	send := make(chan []byte)
//	websocketConn.Run(conn, send, websocketConn.DefaultConfig(), func(bytes []byte) {
//
//	})
//})

// ReadPump pumps messages from the WebSocket connection to the deamon.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (x *webSocketWorker) ReadPump() (err error) {

	for {
		messageType, message, err := x.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				println("WebSocket error 4:", err)
			}
			break
		}
		switch messageType {
		case websocket.TextMessage | websocket.BinaryMessage:
			x.f(message)
		case websocket.CloseMessage:
			break
		}
	}
	return
}

// WritePump pumps messages to the websocket connection from the deamon.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (x *webSocketWorker) WritePump() {

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod := (x.pongWait * 9) / 10

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		x.done <- true
	}()
	for {
		select {
		case <-x.interrupt:
			return
		case message, ok := <-x.send:
			x.conn.SetWriteDeadline(time.Now().Add(x.writeWait))
			if !ok {
				// The server closed the channel.
				x.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := x.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				println("WebSocket: error 1:", err)
				return
			}
			_,err = w.Write(message)
			if err != nil {
				println("WebSocket: error 2:", err)
				return
			}
		case <-ticker.C:
			x.conn.SetWriteDeadline(time.Now().Add(x.writeWait))
			if err := x.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				println("WebSocket: error 3:", err)
				return
			}
		}
	}
}
