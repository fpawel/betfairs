package football

import (
	"sync"
	"log"

)

type GamesReader struct {
	muConsumers   sync.RWMutex
	consumers     []chan resultReadGames

}

type resultReadGames struct {
	games []Game
	error error
}




func (x *GamesReader) read() {
	var r resultReadGames
	r.games, r.error = FetchGames()
	x.muConsumers.Lock()
	consumers := make([]chan resultReadGames, len(x.consumers))
	copy(consumers, x.consumers )
	x.consumers = nil
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

