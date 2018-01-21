package football2

import (
	"time"
	"heroku.com/betfairs/aping/listMarketCatalogue"
	"heroku.com/betfairs/aping/listMarketBook"
	"heroku.com/betfairs/aping"
	"fmt"
	"heroku.com/betfairs/football"
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
