package football

import (
	"sync"
	"log"

)

type SyncReader struct {
	muConsumers   sync.RWMutex
	consumers     []chan gamesResult

}

type gamesResult struct {
	games []Game
	error error
}




func (x *SyncReader) read() {
	var r gamesResult
	r.games, r.error = FetchGames()
	x.muConsumers.Lock()
	consumers := make([]chan gamesResult, len(x.consumers))
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

func (x *SyncReader) Read()  ([]Game, error) {

	ch := make(chan gamesResult)
	x.muConsumers.Lock()
	if x.consumers == nil {
		go x.read()
	}
	x.consumers = append(x.consumers, ch)
	x.muConsumers.Unlock()
	r := <-ch
	return r.games, r.error
}

