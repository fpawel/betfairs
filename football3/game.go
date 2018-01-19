package football3

import (
	"time"

	"heroku.com/betfairs/aping/listMarketBook"
	"heroku.com/betfairs/aping/listMarketCatalogue"
	"heroku.com/betfairs/aping"
)

type Game struct {
	ID int
	OpenDate time.Time `json:"open_date"`
	ScoreHome int    `json:"score_home"`
	ScoreAway int    `json:"score_away"`
	Time      string `json:"time"`
	Markets aping.MarketBooks `json:"markets"`
}

func (x *Game) Read(marketCatalogueReader *listMarketCatalogue.Reader, marketBookReader  *listMarketBook.Reader) error {
	mc,err := marketCatalogueReader.Read(x.ID)
	if err != nil {
		return err
	}
	x.OpenDate = mc[0].Event.OpenDate
	for _,ids := range mc.Take40MarketIDs() {
		mb, err := marketBookReader.Read(ids, time.Second)
		if err == nil {
			x.Markets = append(x.Markets, mb...)
		}
	}
	return nil
}
