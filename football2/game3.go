package football2

import (
	"time"
	"heroku.com/betfairs/aping/listMarketCatalogue"
	"heroku.com/betfairs/aping/listMarketBook"
	"heroku.com/betfairs/aping"
	"fmt"
	"heroku.com/betfairs/football"
	"strconv"
	"log"
)

type Game3 struct {
	Game football.Game
	Event aping.Event
	Competition aping.Competition
	Markets []NamedMarketBook `json:"markets"`
}

type NamedMarketBook struct {

	aping.MarketBookInfo

	//  The name of the market
	Name string `json:"marketName,omitempty"`

	//Information about the runners (selections) in the market.
	Runners []NamedRunner `json:"runners,omitempty"`
}


type NamedRunner struct {

	aping.Runner

	// The name of the runner
	Name string `json:"runnerName,omitempty"`
}

func (x NamedRunner) Runner4() (r Runner4){
	r.ID = x.ID
	r.Name = x.Name
	r.LastPriceTraded = x.LastPriceTraded
	r.Status = x.Status
	r.AvailableToLay = x.ExchangePrices.AvailableToLay.Dub()
	r.AvailableToBack = x.ExchangePrices.AvailableToBack.Dub()
	return
}

func (x NamedMarketBook) Market4() (r Market4){
	var err error
	r.ID, err =  strconv.Atoi( string(x.ID[2:len(x.ID)]) )
	if err != nil {
		log.Fatalf("%s: %v", x.ID, err)
	}
	r.Name = x.Name
	r.TotalAvailable = x.TotalAvailable
	r.TotalMatched = x.TotalMatched
	for _,runner := range  x.Runners {
		r.Runners = append(r.Runners, runner.Runner4())
	}
	return
}

func (x Game3) Game4() (Game4) {
	r := Game4{
		ID:  x.Game.ID,
		Home: x.Game.Home,
		Away: x.Game.Away,
		ScoreHome: x.Game.ScoreHome,
		ScoreAway: x.Game.ScoreAway,
		CountryCode:x.Event.CountryCode,
		CompetitionName:x.Competition.Name,
		OpenDate:x.Event.OpenDate,
	}

	var err error
	r.CompetitionID, err =  strconv.Atoi( x.Competition.ID )
	if err != nil {
		log.Fatalf("%s: %v", x.Competition.ID, err)
	}

	for _,m := range x.Markets{
		r.Markets = append(r.Markets, m.Market4())
	}
	return r
}


func (x *Game3) Read(marketCatalogueReader *listMarketCatalogue.Reader, marketBookReader  *listMarketBook.Reader) error {
	mc,err := marketCatalogueReader.Read(x.Game.ID)
	if err != nil {
		return err
	}
	x.Event = mc[0].Event
	x.Competition = mc[0].Competition

	for _,ids := range mc.Take40MarketIDs() {
		mb, err := marketBookReader.Read(ids, time.Second)
		if err != nil {
			continue
		}

		for _,marketBook := range mb{
			marketCatalogueFound := false
			for _,marketCatalogue := range mc{
				if marketCatalogue.ID == marketBook.ID {
					marketCatalogueFound = true
					marketBook2 := NamedMarketBook{
						MarketBookInfo:marketBook.MarketBookInfo,
						Name:marketCatalogue.Name,
					}
					for _,runner := range marketBook.Runners{
						runnerFound := false
						for _,runnerCatalogue := range marketCatalogue.Runners{
							if runnerCatalogue.ID == runner.ID{
								runner2 := NamedRunner{
									Runner:runner,
									Name:runnerCatalogue.Name,
								}
								marketBook2.Runners = append(marketBook2.Runners, runner2)
								runnerFound = true
								break
							}
						}
						if !runnerFound {
							return fmt.Errorf("%s %s: market runner not found", marketBook.ID, runner.ID)
						}
					}
					x.Markets = append(x.Markets, marketBook2)
					break
				}
			}
			if !marketCatalogueFound {
				return fmt.Errorf("%s: market not found", marketBook.ID)
			}
		}
	}
	return nil
}
