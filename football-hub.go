package main

import (
	"fmt"
	"github.com/fpawel/betfairs/football/football2"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type footballHub struct {
	peers   []*footballPeer
	mu      sync.Mutex
	games   football2.Games
	betfair BetfairClient
}

type footballPeer struct {
	ws   *websocket.Conn
	init bool
}

func newFootballHub(betfair BetfairClient) *footballHub {
	x := &footballHub{
		betfair: betfair,
	}
	go func() {
		for {
			games, err := x.betfair.ReadFootballGames2()
			if err != nil {
				log.Println("Read football failed:", err)
			} else {
				x.broadcast(games)
			}

		}
	}()
	return x
}

func (x *footballHub) add(ws *websocket.Conn) {
	x.mu.Lock()
	x.peers = append(x.peers, &footballPeer{ws: ws})
	x.mu.Unlock()
}

func (x *footballHub) drop(c *footballPeer) {
	x.mu.Lock()
	for i := range x.peers {
		if x.peers[i] == c {
			x.peers[i] = x.peers[len(x.peers)-1]
			x.peers = x.peers[:len(x.peers)-1]
			break
		}
	}
	x.mu.Unlock()
}

func (x *footballHub) Peers() (xs []*footballPeer) {
	x.mu.Lock()
	xs = append([]*footballPeer{}, x.peers...)
	x.mu.Unlock()
	return
}

func (x *footballHub) broadcast(games football2.Games) {
	changes1 := x.games.Changes(games)
	var changes0 *football2.GamesChanges
	wg := new(sync.WaitGroup)
	for _, c := range x.Peers() {
		changes := &changes1
		if !c.init {
			if changes0 == nil {
				tmp := football2.Games{}.Changes(games)
				changes0 = &tmp
			}
			changes = changes0
			c.init = true
		}
		if changes.Empty() {
			continue
		}

		if err := c.ws.WriteJSON(changes); err != nil {
			fmt.Println("WriteJSON: ", err)
			x.drop(c)
		}
	}
	wg.Wait()
	x.games = games
}
