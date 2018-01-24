package football

import (
	"sync"
	"log"

)

type GamesReader struct {
	muConsumers   sync.RWMutex
	consumers     []chan resultReadGames
	teams map[int] teams

}

type teams struct {
	home,away string
}

type resultReadGames struct {
	games []Game
	error error
}


func (x *GamesReader) TeamsByID(id int) (string, string, bool){
	x.muConsumers.RLock()
	t,ok := x.teams[id]
	x.muConsumers.RUnlock()
	if ok {
		return t.home, t.away, true
	}
	return "", "", false
}

func (x *GamesReader) read() {
	var r resultReadGames

	r.games, r.error = FetchGames()

	x.muConsumers.Lock()
	consumers := make([]chan resultReadGames, len(x.consumers))
	copy(consumers, x.consumers )
	x.consumers = nil
	if r.error == nil {
		if x.teams == nil {
			x.teams = make(map[int]teams)
		}
		for _,game := range r.games{
			x.teams[game.ID] = teams{game.Home, game.Away}
		}
	}
	x.muConsumers.Unlock()

	if len(consumers) == 0 {
		log.Fatalln("consumers list is empty")
	}

	for _, consumer := range consumers {
		consumer <- r
	}
	//fmt.Println(len(consumers), "consumers")
}

func (x *GamesReader) Read()  ([]Game, error) {

	ch := make(chan resultReadGames)
	x.muConsumers.Lock()
	if x.consumers == nil {
		go x.read()
	}
	x.consumers = append(x.consumers, ch)
	x.muConsumers.Unlock()
	r := <-ch
	return r.games, r.error
}

