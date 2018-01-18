package aping

import (
	"time"
	"log"
)

///The dynamic data about runners in a market
type Runner struct {

	// The unique id of the runner (selection)
	ID RunnerID `json:"selectionId"`

	// The status of the selection (i.e., ACTIVE, REMOVED, WINNER, LOSER, HIDDEN)
	// Runner status information is available for 90 days following market settlement.
	Status string `json:"status,omitempty"`

	// The adjustment factor applied if the selection is removed
	AdjustmentFactor float64 `json:"adjustmentFactor,omitempty"`

	// The price of the most recent bet matched on this selection
	LastPriceTraded float64 `json:"lastPriceTraded,omitempty"`

	// The total amount matched on this runner
	TotalMatched float64 `json:"totalMatched,omitempty"`

	// If date and time the runner was removed
	RemovalDate time.Time `json:"removalDate,omitempty"`

	// The Exchange prices available for this runner
	ExchangePrices ExchangePrices `json:"ex,omitempty"`

	// The BSP related prices for this runner
	SP  *StartingPrices `json:"sp,omitempty"`

	//The sort priority of this runner
	//SortPriority int `json:"sortPriority,omitempty"`

	// The handicap.  Enter the specific handicap value (returned by RUNNER in listMaketBook)
	// if the market is an Asian handicap market.
	//Handicap float64 `json:"handicap,omitempty"`

	// List of orders in the market
	//orders : Order list

	// List of matches (i.e, orders that have been fully or partially executed)
	//matches : Match list
}



func (x Runner) PriceSize(side Side) []PriceSize {
	switch side {
	case Back:
		return x.ExchangePrices.AvailableToBack
	case Lay:
		return  x.ExchangePrices.AvailableToLay
	default:
		log.Fatalf("side must be Back or Lay, but it is %v", side)
	}
	return nil
}



func (x Runner) Dub() (dub Runner)  {
	dub = x
	dub.ExchangePrices = x.ExchangePrices.Dub()
	return
}


type StartingPrices struct {

	//What the starting price would be if the market was reconciled now taking into account the SP bets as well as unmatched exchange bets on the same selection in the exchange. This data is cached and update every 60 seconds. Please note: Type Double may contain numbers, INF, -INF, and NaN.
	NearPrice	float64 `json:"nearPrice,omitempty"`


	//What the starting price would be if the market was reconciled now taking into account only the currently place SP bets. The Far Price is not as complicated but not as accurate and only accounts for money on the exchange at SP. This data is cached and updated every 60 seconds. Please note: Type Double may contain numbers, INF, -INF, and NaN.
	FarPrice float64 `json:"farPrice,omitempty"`

	//The total amount of back bets matched at the actual Betfair Starting Price.
	BackStakeTaken []PriceSize `json:"backStakeTaken,omitempty"`

	//The lay amount matched at the actual Betfair Starting Price.
	LayLiabilityTaken []PriceSize `json:"layLiabilityTaken,omitempty"`

	// The final BSP price for this runner. Only available for a BSP market that has been reconciled.
	ActualSP float64 `json:"actualSP,omitempty"`


}