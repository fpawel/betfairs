package aping

import "time"

type MarketCatalogue struct {

	//The unique identifier for the market. MarketId's are prefixed with '1.' or '2.' 1. = UK Exchange 2. = AUS Exchange.
	ID string `json:"marketId"`

	//  The name of the market
	Name string `json:"marketName,omitempty"`

	//The total amount matched
	TotalMatched float64 `json:"totalMatched,omitempty"`

	// The time this market starts at, only returned when the MARKET_START_TIME enum is passed in the marketProjections
	MarketStartTime time.Time `json:"marketStartTime,omitempty"`

	// Details about the market
	Description string `json:"description,omitempty"`

	//The runners (selections) contained in the market
	Runners []RunnerCatalogue `json:"runners,omitempty"`

	// The competition the market is contained within. Usually only applies to Football competitions
	Competition Competition `json:"competition,omitempty"`

	// The Event Type the market is contained within
	EventType EventType `json:"eventType,omitempty"`

	// The Event the market is contained within
	Event Event `json:"event,omitempty"`
}

func (x MarketCatalogue) Runner(id int) (r RunnerCatalogue, ok bool)  {
	for i := range  x.Runners{
		if x.Runners[i].ID == id {
			r = x.Runners[i]
			ok = true
			break
		}
	}
	return
}

func (x MarketCatalogue) Dub() (dub MarketCatalogue)  {
	dub = x
	dub.Runners = make([]RunnerCatalogue, len(x.Runners))
	for i := range  x.Runners{
		dub.Runners[i] = x.Runners[i]
	}
	return
}