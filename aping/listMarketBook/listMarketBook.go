package listMarketBook

import (
	"sync"
	"time"
	"log"
	"heroku.com/betfairs/aping"

)

type Reader struct {
	muConsumers   sync.RWMutex
	consumers     map[aping.MarketID][] chan<- resultRead
	muCache      sync.RWMutex
	cache        map[aping.MarketID]cachedItem
	apingSession *aping.Session
}



func New(apingSession *aping.Session) (x *Reader) {
	return &Reader{
		apingSession:apingSession,
		consumers : make(map[aping.MarketID][] chan<- resultRead),
		cache : make(map[aping.MarketID]cachedItem),
	}
}

type cachedItem struct {
	marketBook aping.MarketBook
	time       time.Time
}

type resultRead struct {
	Markets aping.MarketBooks
	Error error
}

func (x *Reader) read(marketIDs []aping.MarketID) {

	if len(marketIDs) == 0{
		log.Fatal("marketIDs must be not nil")
	}

	var r resultRead
	r.Markets, r.Error = x.apingSession.ListMarketBook(marketIDs)
	if r.Error == nil {
		x.muCache.Lock()
		// добавить считаное в кеш
		for _,m := range r.Markets{
			x.cache[m.ID] = cachedItem{m, time.Now()}
		}

		x.muCache.Unlock()
	}

	x.muConsumers.Lock()
	for _,marketID := range marketIDs {
		consumers, ok := x.consumers[marketID]
		if !ok || len(consumers) == 0 {
			log.Fatalln("consumers list is empty")
		}
		for _, consumer := range consumers {
			consumer <- r
		}
		delete(x.consumers, marketID)
	}
	x.muConsumers.Unlock()
}

func (x *Reader) doRead( marketIDs []aping.MarketID) (result aping.MarketBooks, err error) {
	if len(marketIDs) == 0{
		log.Fatal("marketIDs must be not nil")
	}
	var newConsumers []chan resultRead
	x.muConsumers.Lock()

	var readMarketIDs []aping.MarketID

	for _,marketID := range marketIDs {
		ch := make(chan resultRead)
		newConsumers = append(newConsumers, ch)

		existedConsumers, updating := x.consumers[marketID]
		x.consumers[marketID] = append(existedConsumers, ch)
		if !updating {
			readMarketIDs = append(readMarketIDs, marketID)
		}
	}
	x.muConsumers.Unlock()
	if len(readMarketIDs) > 0 {
		go x.read(readMarketIDs)
	}

	for _, ch := range newConsumers{
		r := <- ch
		result = append(result, r.Markets ...)
		err = r.Error
	}
	return
}

func (x *Reader) Read(marketIDs []aping.MarketID, t time.Duration)  (result aping.MarketBooks, err error){

	result = make( aping.MarketBooks, len(marketIDs) )
	posInResult := make(map[aping.MarketID]int)

	x.muCache.Lock()

	// вычистить из кеша тухляк
	for k,c := range x.cache {
		if time.Since(c.time) > time.Minute * 5 {
			delete(x.cache, k)
		}
	}

	var readMarketIDs []aping.MarketID

	for i,marketID := range marketIDs {
		posInResult[marketID] = i
		cached, found := x.cache[marketID]
		if found && time.Since(cached.time) < t {
			result[i] = cached.marketBook.Dub()
		} else {
			readMarketIDs = append(readMarketIDs, marketID)
		}
	}

	x.muCache.Unlock()

	if len(readMarketIDs) > 0 {
		var xs aping.MarketBooks
		xs,err = x.doRead(readMarketIDs)
		if err == nil {
			for _,marketBook := range xs {
				n,ok := posInResult[marketBook.ID]
				if !ok {
					log.Printf("market %s not found", marketBook.ID)
					continue
				}
				if n<0 || n>= len(marketIDs) {
					log.Fatal("out of range")
				}
				result[n] = marketBook.Dub()

				for i,id := range readMarketIDs{
					if id == marketBook.ID {
						readMarketIDs[i] = readMarketIDs[len(readMarketIDs)-1]
						readMarketIDs = readMarketIDs[:len(readMarketIDs)-1]
					}
				}
			}
			if len(readMarketIDs) > 0 {
				log.Fatal("readMarketIDs must be empty at the end")
			}
		} else {
			result = nil
		}
	}
	return
}


func (x *Reader) Get(marketID aping.MarketID)  (result aping.MarketBook, ok bool){

	x.muCache.Lock()
	// вычистить из кеша тухляк
	for k,c := range x.cache {
		if time.Since(c.time) > time.Minute * 5 {
			delete(x.cache, k)
		}
	}

	cached, found := x.cache[marketID]
	if found  {
		result = cached.marketBook.Dub()
		ok = true
	} else {
		go func() {
			x.Read([]aping.MarketID{marketID}, time.Second)
		}()
	}
	x.muCache.Unlock()

	return
}
