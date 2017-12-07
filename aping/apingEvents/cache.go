package apingEvents

import (
	"log"
	"sync"
	"time"
	"heroku.com/betfairs/aping"
)

type Cache struct {
	muConsumers   sync.RWMutex
	consumers     map[int][]consumer
	muCache      sync.RWMutex
	cache        map[int]cachedItem
	apingSession *aping.Session

	updating bool
}

func NewCache(apingSession *aping.Session) (x *Cache) {
	return &Cache{
		apingSession:apingSession,
		consumers : make(map[int][] consumer),
		cache : make(map[int]cachedItem),
	}
}

type cachedItem struct {
	event *aping.Event
	time  time.Time
}

type consumer struct {
	handler func (*aping.Event,error)
	wait chan bool
}

func (x *Cache) readEvent(eventID int) {

	if x.updating {
		log.Fatal("invariant failed: Cache updating")
	}
	x.updating = true
	event, err := x.apingSession.ReadEvent(eventID)
	x.updating = false

	if err == nil {
		x.muCache.Lock()
		// добавить считаное в кеш
		x.cache[eventID] = cachedItem{event, time.Now()}
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
		consumer.handler(event, err)
		consumer.wait <- true
	}
	//fmt.Println(len(consumers), "ready now")
}

func (x *Cache) ViewEvent(eventID int, handler func (*aping.Event,error))  <-chan bool{

	wait := make(chan bool)

	x.muCache.RLock()
	if cached, found := x.cache[eventID]; found  {
		handler(cached.event, nil)
		x.muCache.RUnlock()
		go func () {
			wait <- true
		}()
		return wait
	}
	x.muCache.RUnlock()
	consumer := consumer{
		wait: make(chan bool),
		handler:handler,
	}

	x.muConsumers.Lock()
	defer x.muConsumers.Unlock()

	consumers, updating := x.consumers[eventID]
	if updating != x.updating {
		log.Fatal("invariant failed: Cache updating")
	}

	x.consumers[eventID] = append(consumers, consumer)

	if !updating {
		go x.readEvent(eventID)
	}

	return consumer.wait
}

func (x *Cache) ClearEvent(eventID int) {
	x.muCache.Lock()
	defer x.muCache.Unlock()
	delete(x.cache, eventID)
}

