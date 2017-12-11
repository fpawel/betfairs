package football2

import (
	"heroku.com/betfairs/football"

	"heroku.com/betfairs/aping/listMarketCatalogue"
	"fmt"
	"heroku.com/betfairs/aping/listMarketBook"
	"heroku.com/betfairs/countries"
	"time"
	"heroku.com/betfairs/aping"
)

type Game struct {
	football.Game
	Competition string `json:"competition"`
	Country string `json:"country"`
	MainPrices [6]float64 `json:"main_prices"`
	TotalMatched float64 `json:"total_matched"`
	TotalAvailable float64 `json:"total_available"`
	Error error `json:"error"`
}



func (x *Game) Read(marketCatalogueReader *listMarketCatalogue.Reader, marketBookReader  *listMarketBook.Reader)  {
	var mc aping.MarketCatalogues
	mc,x.Error = marketCatalogueReader.Read(x.ID)
	if x.Error != nil {
		return
	}
	x.Competition = mc[0].Competition.Name
	c := countries.ByAlpha2(mc[0].Event.CountryCode)
	if c != nil {
		x.Country = c.Name
	} else {
		x.Country = mc[0].Event.CountryCode
	}
	mainMarket,ok := mc.MainMarket()
	if !ok {
		x.Error = fmt.Errorf("рынок ставок на результат не найден")
		return
	}

	t := time.Second
	if !x.InPlay {
		t = time.Minute
	}

	var mb aping.MarketBooks
	mb, x.Error = marketBookReader.Read([]string{mainMarket.ID}, t)
	if x.Error != nil {
		return
	}
	x.MainPrices = mb[0].Prices6()
	x.TotalMatched = mb[0].TotalMatched
	x.TotalAvailable = mb[0].TotalAvailable
	x.Error = nil
}
