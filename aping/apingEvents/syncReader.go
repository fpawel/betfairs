package apingEvents

import (
	"log"
	"sync"
	"time"
	"heroku.com/betfairs/aping"
)

type SyncReader struct {
	muConsumers   sync.RWMutex
	consumers     map[int][] chan<- readEventResult
	muCache      sync.RWMutex
	cache        map[int]cachedItem
	apingSession *aping.Session
}

func NewSyncReader(apingSession *aping.Session) (x *SyncReader) {
	return &SyncReader{
		apingSession:apingSession,
		consumers : make(map[int][] chan<- readEventResult),
		cache : make(map[int]cachedItem),
	}
}

type cachedItem struct {
	event *aping.Event
	time  time.Time
}

type readEventResult struct {
	Event *aping.Event
	Error error
}

func (x *SyncReader) readEvent(eventID int) {

	var r readEventResult
	r.Event, r.Error = x.apingSession.ReadEvent(eventID)


	if r.Error == nil {
		x.muCache.Lock()
		// добавить считаное в кеш
		x.cache[eventID] = cachedItem{r.Event, time.Now()}
		// вычистить из кеша тухляк
		for k,c := range x.cache {
			if time.Since(c.time) > time.Hour {
				delete(x.cache, k)
			}
		}
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

func (x *SyncReader) doReadEvent(eventID int) (*aping.Event, error) {
	consumer := make(chan readEventResult)

	x.muConsumers.Lock()
	consumers, updating := x.consumers[eventID]
	x.consumers[eventID] = append(consumers, consumer)
	x.muConsumers.Unlock()

	if !updating {
		go x.readEvent(eventID)
	}
	r := <- consumer
	return r.Event, r.Error
}

func (x *SyncReader) ReadEvent(eventID int)  (*aping.Event, error){

	x.muCache.RLock()
	cached, found := x.cache[eventID]
	x.muCache.RUnlock()

	if found  {
		return cached.event, nil
	}
	return x.doReadEvent(eventID)
}

func (x *SyncReader) Event(eventID int)  *aping.Event {

	x.muCache.RLock()
	cached, found := x.cache[eventID]
	x.muCache.RUnlock()

	if found  {
		return cached.event
	} else {
		go func() {
			x.doReadEvent(eventID)
		}()
		return nil
	}
}

func (x *SyncReader) ClearEvent(eventID int) {
	x.muCache.Lock()
	defer x.muCache.Unlock()
	delete(x.cache, eventID)
}

