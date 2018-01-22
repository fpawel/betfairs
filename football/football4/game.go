package football4

import (
	"time"

	"heroku.com/betfairs/aping"
	"strconv"
	"log"
	"heroku.com/betfairs/football"
	"heroku.com/betfairs/aping/listMarketCatalogue"
	"heroku.com/betfairs/aping/listMarketBook"
)

type Game struct {
	ID              int       `json:"id"`
	OpenDate        time.Time `json:"openDate"`
	ScoreHome       int       `json:"score_home"`
	ScoreAway       int       `json:"score_away"`
	Minute          int       `json:"minute"`
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
	r.ID, err =  strconv.Atoi( string(x.ID[2:len(x.ID)]) )
	if err != nil {
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

func newGame(x football.Game, openDate time.Time, markets aping.MarketBooks) Game {
	r := Game{
		ID:  x.ID,
		ScoreHome: x.ScoreHome,
		ScoreAway: x.ScoreAway,
		OpenDate:openDate,
		Minute:-1,
	}
	if minute,err := x.Minute(); err == nil {
		r.Minute = minute
	}
	for _,m := range markets{
		r.Markets = append(r.Markets, newMarket(m))
	}
	return r
}

func ReadGame(game football.Game,
	marketCatalogueReader *listMarketCatalogue.Reader,
	marketBookReader  *listMarketBook.Reader) (Game, error) {
	var marketCatalogues aping.MarketCatalogues
	marketCatalogues,err := marketCatalogueReader.Read(game.ID)
	if err != nil {
		var tmp Game
		return tmp,err
	}

	var markets aping.MarketBooks
	for _,ids := range marketCatalogues.Take40MarketIDs() {
		mb, err := marketBookReader.Read(ids, time.Second)
		if err != nil {
			continue
		}
		markets = append(markets, mb...)
	}
	return newGame(game, marketCatalogues[0].Event.OpenDate, markets), nil
}


