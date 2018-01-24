package event2

import (
	"github.com/fpawel/betfairs/aping"
	"log"
	"strconv"
	"time"
)

type Event struct {
	ID          int       `json:"id"`
	OpenDate    time.Time `json:"openDate"`
	Competition aping.Competition
	Home        string   `json:"home"`
	Away        string   `json:"away"`
	Markets     []Market `json:"markets"`
	CountryCode string `json:"countryCode"`

}

type Market struct {
	ID      int                     `json:"marketId"`
	Name    string                  `json:"marketName,omitempty"`
	Runners []aping.RunnerCatalogue `json:"runners,omitempty"`
}

func NewEvent(eventID int, xs aping.MarketCatalogues, home, away string) Event {
	if len(xs) == 0 {
		log.Fatal(aping.ErrorNoMarkets)
	}
	ev := xs[0].Event
	n,err := strconv.Atoi(ev.ID)
	if err != nil {
		log.Fatal(err)
	}
	if n != eventID {
		log.Fatalf("%d != %d", n, eventID)
	}


	return Event{
		ID: eventID,
		Home:        home,
		Away:        away,
		OpenDate:    ev.OpenDate,
		Competition: xs[0].Competition,
		Markets:     markets(xs),
		CountryCode: ev.CountryCode,
	}
}

func markets(xs aping.MarketCatalogues) (r []Market) {

	for _, m := range xs {
		var err error
		a := Market{
			Name:    m.Name,
			Runners: m.Dub().Runners,
		}
		a.ID, err = strconv.Atoi(string(m.ID[2:len(m.ID)]))
		if err != nil {
			log.Fatalf("%s: %v", m.ID, err)
		}
		r = append(r, a)
	}
	return
}
