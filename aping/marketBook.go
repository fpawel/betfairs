package aping

import "time"

type MarketBook struct {

	//The unique identifier for the market. MarketId's are prefixed with '1.' or '2.' 1. = UK Exchange 2. = AUS Exchange.
	ID MarketID `json:"marketId"`

	//True if the data returned by listMarketBook will be delayed.
	// The data may be delayed because you are not logged in with a funded account or you are using an Application Key
	// that does not allow up to date data.
	IsMarketDataDelayed bool `json:"isMarketDataDelayed,omitempty"`

	//The status of the market, for example ACTIVE, SUSPENDED, SETTLED, etc.
	Status string `json:"status,omitempty"`

	// The number of seconds an order is held until it is submitted into the market.
	// Orders are usually delayed when the market is in-play
	BetDelay int `json:"betDelay,omitempty"`

	// True if the market starting price has been reconciled
	BspReconciled bool `json:"bspReconciled,omitempty"`

	//If false, runners may be added to the market
	Complete bool `json:"complete,omitempty"`

	// True if the market is currently in play
	Inplay bool `json:"inplay,omitempty"`

	//The number of selections that could be settled as winners
	NumberOfWinners int `json:"numberOfWinners,omitempty"`

	//The number of runners in the market
	NumberOfRunners int `json:"numberOfRunners,omitempty"`

	//The number of runners that are currently active. An active runner is a selection available for betting
	NumberOfActiveRunners int `json:"numberOfActiveRunners,omitempty"`

	//The most recent time an order was executed
	LastMatchTime time.Time `json:"lastMatchTime,omitempty"`

	//The total amount matched
	TotalMatched float64 `json:"totalMatched,omitempty"`

	//The total amount of orders that remain unmatched
	TotalAvailable float64 `json:"totalAvailable,omitempty"`

	//True if cross matching is enabled for this market.
	CrossMatching bool `json:"crossMatching,omitempty"`

	//True if runners in the market can be voided
	RunnersVoidable bool `json:"runnersVoidable,omitempty"`

	//The version of the market. The version increments whenever the market status changes,
	// for example, turning in-play, or suspended when a goal score.
	Version int `json:"version,omitempty"`

	//Information about the runners (selections) in the market.
	Runners []Runner `json:"runners,omitempty"`


}

func (x MarketBook) Prices6() (r [6]float64){

	if len(x.Runners) != 3 {
		return
	}
	r[0] = x.Runners[0].ExchangePrices.Back()
	r[1] = x.Runners[0].ExchangePrices.Lay()
	r[2] = x.Runners[1].ExchangePrices.Back()
	r[3] = x.Runners[1].ExchangePrices.Lay()
	r[4] = x.Runners[2].ExchangePrices.Back()
	r[5] = x.Runners[2].ExchangePrices.Lay()
	return
}

func (x MarketBook) Dub() (dub MarketBook)  {
	dub = x
	//if x.LastMatchTime != nil {
	//	*dub.LastMatchTime = *x.LastMatchTime
	//}

	dub.Runners = make([]Runner, len(x.Runners))
	for i := range  x.Runners{
		dub.Runners[i] = x.Runners[i].Dub()
	}
	return
}

type PriceProjection struct {
	PriceData  []string `json:"priceData"`
	Virtualise bool     `json:"virtualise"`
}
