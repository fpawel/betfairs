package event2

import "heroku.com/betfairs/aping"

type Event struct {
	Score, Time string

	Odds map[MarketID] map [RunnerID]  *[2]float64
}


type Market struct {
	ID aping.MarketID

}
