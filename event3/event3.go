package event3

import (
	"time"

	"github.com/fpawel/betfairs/aping"
	"strconv"
	"log"
	"github.com/fpawel/betfairs/aping/listMarketCatalogue"
	"github.com/fpawel/betfairs/aping/listMarketBook"
)

type Event struct {
	ID              int       `json:"id"`
	Markets []Market          `json:"markets"`
}

type Market struct {
	ID int                 `json:"marketId"`
	TotalMatched float64   `json:"totalMatched,omitempty"`
	TotalAvailable float64 `json:"totalAvailable,omitempty"`
	Runners []Runner       `json:"runners,omitempty"`
}

type Runner struct {
	ID              aping.RunnerID   `json:"selectionId"`
	Status          string           `json:"status,omitempty"` // ACTIVE, REMOVED, WINNER, LOSER, HIDDEN
	LastPriceTraded float64          `json:"lastPriceTraded,omitempty"`
	AvailableToBack aping.PriceSizes `json:"availableToBack,omitempty"`
	AvailableToLay  aping.PriceSizes `json:"availableToLay,omitempty"`
}

func newRunner(x aping.Runner) (r Runner){
	r.ID = x.ID
	r.LastPriceTraded = x.LastPriceTraded
	r.Status = x.Status
	r.AvailableToLay = x.ExchangePrices.AvailableToLay.Dub()
	r.AvailableToBack = x.ExchangePrices.AvailableToBack.Dub()
	return
}

func newMarket(x aping.MarketBook) (r Market){
	var err error
	r.ID = x.ID.Int()
	if r.ID == 0 {
		log.Fatalf("%s: %v", x.ID, err)
	}
	r.TotalAvailable = x.TotalAvailable
	r.TotalMatched = x.TotalMatched
	for _,runner := range  x.Runners {
		ex := runner.ExchangePrices
		if len(ex.AvailableToBack)==0  && len(ex.AvailableToLay ) == 0 && runner.Status == "ACTIVE"{
			continue
		}
		r.Runners = append(r.Runners, newRunner(runner))
	}
	return
}

func newGame(eventID int, markets aping.MarketBooks) Event {
	r := Event{
		ID:  eventID,
	}
	for _,m := range markets{
		r.Markets = append(r.Markets, newMarket(m))
	}
	return r
}

func ReadEvent(eventID int,	mcr *listMarketCatalogue.Reader, mbr *listMarketBook.Reader) (Event, error) {
	var marketCatalogues aping.MarketCatalogues
	marketCatalogues,err := mcr.Read(eventID)
	if err != nil {
		var tmp Event
		return tmp,err
	}

	var markets aping.MarketBooks
	for _,ids := range marketCatalogues.Take40MarketIDs() {
		mb, err := mbr.Read(ids, time.Second)
		if err != nil {
			continue
		}
		markets = append(markets, mb...)
	}
	return newGame(eventID, markets), nil
}


