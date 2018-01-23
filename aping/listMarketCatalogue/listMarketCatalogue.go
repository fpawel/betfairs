package listMarketCatalogue

import (
	"log"
	"sync"
	"time"
	"github.com/fpawel/betfairs/aping"
)

type Reader struct {
	muConsumers   sync.RWMutex
	consumers     map[int][] chan<- resultRead
	muCache      sync.RWMutex
	cache        map[int]cachedItem
	apingSession *aping.Session

}

func New(apingSession *aping.Session) (x *Reader) {
	return &Reader{
		apingSession:apingSession,
		consumers : make(map[int][] chan<- resultRead),
		cache : make(map[int]cachedItem),
	}
}

type cachedItem struct {
	markets aping.MarketCatalogues
	time    time.Time
}

type resultRead struct {
	Markets aping.MarketCatalogues
	Error   error
}

func (x *Reader) read(eventID int) {

	var r resultRead
	r.Markets, r.Error = x.apingSession.ListMarketCatalogue(eventID)


	if r.Error == nil {
		x.muCache.Lock()
		// добавить считаное в кеш
		x.cache[eventID] = cachedItem{r.Markets, time.Now()}
		x.muCache.Unlock()
	}

	x.muConsumers.Lock()
	consumers, ok := x.consumers[eventID]
	delete(x.consumers, eventID)
	x.muConsumers.Unlock()

	if !ok || len(consumers) == 0 {
		log.Fatalln("consumers list is empty")
	}

	for _, consumer := range consumers {
		consumer <- r
	}
	//fmt.Println(len(consumers), "ready now")
}

func (x *Reader) update(eventID int) ([]aping.MarketCatalogue, error) {
	consumer := make(chan resultRead)

	x.muConsumers.Lock()
	consumers, updating := x.consumers[eventID]
	x.consumers[eventID] = append(consumers, consumer)
	x.muConsumers.Unlock()

	if !updating {
		go x.read(eventID)
	}
	r := <- consumer

	return r.Markets.Dub(), r.Error
}

func (x *Reader) Read(eventID int)  (aping.MarketCatalogues, error){

	x.muCache.Lock()
	// вычистить из кеша тухляк
	for k,c := range x.cache {
		if time.Since(c.time) > time.Hour {
			delete(x.cache, k)
		}
	}
	cached, found := x.cache[eventID]
	x.muCache.Unlock()

	if found  {
		return cached.markets.Dub() , nil
	}
	return x.update(eventID)
}


func (x *Reader) Get(eventID int)  (aping.MarketCatalogues, bool) {

	var cached cachedItem
	x.muCache.Lock()
	// вычистить из кеша тухляк
	for k,c := range x.cache {
		if time.Since(c.time) > time.Hour {
			delete(x.cache, k)
		}
	}
	cached, ok := x.cache[eventID]
	x.muCache.Unlock()

	if ok  {
		return cached.markets.Dub(), true

	} else {
		go func() {
			x.update(eventID)
		}()
	}
	return nil, false
}



